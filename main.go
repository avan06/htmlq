package main

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/andybalholm/cascadia"
	"github.com/antchfx/htmlquery"
	"github.com/fatih/color"
	"golang.org/x/net/html"
)

var (
	htmlData  string
	results   [][]interface{}
	headerMap = make(map[int]interface{})
	// Create new parser object
	parser = argparse.NewParser("htmlq 1.0.5", "A command-line tool that allows you to query HTML using CSS selectors or XPATH and retrieve the corresponding text content (similar to JavaScript's `document.querySelector(query).textContent`)")

	// Create filePath flag
	filePath = parser.String("f", "file", &argparse.Options{Help: `Enter the relative or absolute path of the HTML file`})
	// Create inputText flag
	inputText = parser.String("t", "text", &argparse.Options{Help: `Enter the HTML text content`})
	// Create inputUrl flag
	inputUrl = parser.String("u", "url", &argparse.Options{Help: `Enter the URL of the HTML`})

	// Create isXPATH flag
	isXPATH = parser.Flag("x", "XPATH", &argparse.Options{Help: `Enable default XPATH query syntax instead of CSS Selectors`, Default: false})
	// Create isSelectorAll flag
	isSelectorAll = parser.Flag("a", "SelectorAll", &argparse.Options{Help: `Enable the SelectorAll mechanism. When there are multiple results for a single query, it will return all the results (similar to querySelectorAll). By default, it is not enabled, and only a single result is returned (similar to querySelector).`, Default: false})
	// Create isResultAsNode flag
	isResultAsNode = parser.Flag("r", "ResultAsNode", &argparse.Options{Help: `Enable using the Node from the previous query result as the current query's root Node`, Default: false})
	// Create isPrintLastResult flag
	isPrintLastResult = parser.Flag("l", "PrintLastResult", &argparse.Options{Help: `Enable printing the content of the last result in the output when using the "#lastresult" syntax in the query`, Default: false})
	// Create isPrintLastResultTemp flag
	isPrintLastResultTemp = parser.Flag("", "PrintLastResultTemp", &argparse.Options{Help: `Enable printing the temporary content of the source data as last result in the output, when using the "#lastresult" syntax in query`, Default: false})
	// Create isPrintQuerysInAllResult flag
	isPrintQuerysInAllResult = parser.Flag("", "PrintQuerysInAllResult", &argparse.Options{Help: `Enable printing of query conditions in results.`, Default: false})
	// Create PrintTextToAllResult flag
	printTextInAllResult = parser.String("", "PrintTextInAllResult", &argparse.Options{Help: `Enter a text content for printing in all results. If the content contains "#{serial number}", the content of the specified serial number's result will be automatically printed in all results following that serial number.`})

	// Create headers flag
	headers = parser.StringList("H", "headers", &argparse.Options{Help: `When the query item is a table or multiple td fields, you can enter corresponding names for each individual field in a single query using the format "#{serial number}:header1Name;header2Name;header3Name;...", where the serial number represents the Nth query starting from zero.`})
	// Create verboses flag
	verboses = parser.FlagCounter("v", "verbose", &argparse.Options{Help: `verbose`})
	// Create querys flag
	querys = parser.StringList("q", "querys", &argparse.Options{Required: true, Help: `Enter a query value to retrieve the Node Text of the target HTML. The default query method is CSS Selectors, but it can be changed to XPATH using "-x" or "--XPATH". Both query methods can be mixed. When a query starts with "C:", it represents CSS Selectors, and when it starts with "X:", it represents XPATH.

	You can enter multiple query values separated by spaces and return all the results.
	
	* When the "--ResultNode" is enabled, you can use the following special query values to change the current position of the HTML Node:
	  - Parent: Move the Node up one level
	  - NextSibling: Move the Node to the next sibling in the same level
	  - PrevSibling: Move the Node to the previous sibling in the same level
	  - FirstChild: Move the Node to the first child in the same level
	  - LastChild: Move the Node to the last child in the same level
	  - reset: Restore the Node to the root node of the original input
	
	* When using the "#lastresult" query syntax, the text content of the previous query will automatically replace "#lastresult".
	
	Example: --ResultNode -q Query1 Parent Query2 NextSibling Query3 LastChild C:Query4 reset X:Query5`})
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			color.Red("%v", r)
		}
	}()
	// Parse input
	isDash := false
	for idx := 1; idx < len(os.Args); idx++ {
		arg := os.Args[idx]
		if strings.HasPrefix(arg, "-") {
			isDash = true
		} else if isDash {
			isDash = false
		} else {
			os.Args = append(os.Args[:idx+1], os.Args[idx:]...)
			os.Args[idx] = "-q"
			idx++
		}
	}
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		panic(parser.Usage(err))
	}

	var htmlDoc *html.Node = nil
	if len(*inputText) > 0 {
		htmlData = *inputText
	} else if len(*filePath) > 0 {
		if (*filePath)[0] == '~' { // Check if the path starts with a tilde (~)
			// Split the path into user part and remaining path part
			pathParts := strings.SplitN(*filePath, string(os.PathSeparator), 2)
			userPart := pathParts[0][1:] // Remove the character before the tilde
			remainingPath := ""
			if len(pathParts) > 1 {
				remainingPath = pathParts[1]
			}
			// Check if it is the current user
			if userPart == "" {
				currentUser, err := user.Current()
				if err != nil {
					panic("Unable to get current user information: " + err.Error())
				}
				// Replace the tilde in the path with the current user's home directory
				*filePath = filepath.Join(currentUser.HomeDir, remainingPath)
			} else {
				// Get the home directory path of the specified user
				userAccount, err := user.Lookup(userPart)
				if err != nil {
					panic("Unable to get user information: " + err.Error())
				}
				// Replace the tilde in the path with the specified user's home directory
				*filePath = filepath.Join(userAccount.HomeDir, remainingPath)
			}
		}
		fileContent, err := os.ReadFile(*filePath)
		if err != nil {
			panic(err)
		}
		htmlData = string(fileContent)
	} else if len(*inputUrl) > 0 {
		htmlDoc, err = htmlquery.LoadURL(*inputUrl)
		if err != nil {
			panic(err)
		}
	} else {
		panic(parser.Usage("At least one of 'text' or 'file' or 'url' needs to be provided."))
	}
	if htmlDoc == nil {
		htmlDoc, err = html.Parse(strings.NewReader(htmlData))
		if err != nil {
			panic(err)
		}
	}

	reg := regexp.MustCompile(`\#(\d+):`)
	for _, head := range *headers {
		headMatchs := reg.FindStringSubmatch(head)
		if len(headMatchs) != 2 {
			panic(parser.Usage(fmt.Sprintf("%s head:%s, %v", "The format of the input headers is invalid.", head, headMatchs)))
		}
		queryID, _ := strconv.Atoi(headMatchs[1])
		head = strings.TrimPrefix(head, headMatchs[0])
		heads := strings.Split(head, ";")
		headerMap[queryID] = heads
	}

	var printText string = ""
	var printTextResultID *int = new(int)
	printTextMatchs := regexp.MustCompile(`\#result(\d+)`).FindStringSubmatch(*printTextInAllResult)
	if (len(printTextMatchs)) == 2 {
		*printTextResultID, _ = strconv.Atoi(printTextMatchs[1])
	}

	currDoc := htmlDoc
	var selected []*html.Node = nil
	for qIdx, qVal := range *querys {
		queryXPATH := *isXPATH
		if strings.ToLower(qVal) == "parent" {
			currDoc = currDoc.Parent
			continue
		} else if strings.ToLower(qVal) == "nextsibling" {
			currDoc = currDoc.NextSibling
			continue
		} else if strings.ToLower(qVal) == "prevsibling" {
			currDoc = currDoc.PrevSibling
			continue
		} else if strings.ToLower(qVal) == "firstchild" {
			currDoc = currDoc.FirstChild
			continue
		} else if strings.ToLower(qVal) == "lastchild" {
			currDoc = currDoc.LastChild
			continue
		} else if strings.ToLower(qVal) == "reset" {
			currDoc = htmlDoc
			continue
		} else if strings.HasPrefix(qVal, "X:") {
			qVal = strings.TrimPrefix(qVal, "X:")
			queryXPATH = true
		} else if strings.HasPrefix(qVal, "C:") {
			qVal = strings.TrimPrefix(qVal, "C:")
			queryXPATH = false
		}

		if strings.Contains(qVal, "#lastresult") && len(results) > 0 {
			selected = nil
			for _, data := range results[len(results)-1] {
				tmpDoc := currDoc
				lastresult := fmt.Sprint(data)
				if value, exists := headerMap[len(results)-1]; exists {
					lastresult = strings.TrimPrefix(lastresult, value.([]string)[0]+":")
				}
				lastresult = strings.TrimSuffix(lastresult, "\n")
				valTmp := strings.ReplaceAll(qVal, "#lastresult", lastresult)
				//fmt.Println("valTmp", valTmp)
				selectedTmp := querySelect(queryXPATH, tmpDoc, valTmp)
				selected = append(selected, selectedTmp...)
			}
		} else {
			selected = querySelect(queryXPATH, currDoc, qVal)
		}

		var result []interface{}
		level := 0
		for idx, ele := range selected {
			text := &bytes.Buffer{}
			collectText(ele, text, level)
			//To correctly read the header names, increment the level by one after processing each td element node.
			if idx+1 < len(selected) && ele.Type == html.ElementNode && ele.Data == "td" {
				if ele.Parent == selected[idx+1].Parent {
					level++
				} else {
					level = 0
				}
			}
			result = append(result, strings.TrimSuffix(text.String(), "\n"))
			if !*isSelectorAll {
				break
			}
		}

		if !*isPrintLastResultTemp && qIdx+1 < len(*querys) && strings.Contains((*querys)[qIdx+1], "#lastresult") {
		} else {
			if *verboses > 1 {
				info := fmt.Sprint(len(results), ": ")
				info += qVal
				fmt.Println()
				color.HiGreen(info)
			}

			if strings.Contains(*printTextInAllResult, "#result") && printTextResultID != nil && len(results) > *printTextResultID {
				printText = fmt.Sprintf("%q", results[*printTextResultID])
				printText = strings.NewReplacer(`" "`, ";", "[", "", "]", "", `"`, "").Replace(printText)
				*printTextInAllResult = strings.ReplaceAll(*printTextInAllResult, fmt.Sprintf("#result%v", *printTextResultID), printText)
			}

			for idx, text := range result {
				if *verboses > 0 {
					color.HiCyan("%d-%d:\n", len(results), idx)
					if *isPrintQuerysInAllResult {
						color.HiMagenta(qVal)
					}
					if len(*printTextInAllResult) > 0 && !strings.Contains(*printTextInAllResult, "#result") {
						color.HiYellow(*printTextInAllResult)
					}
					if *isPrintLastResult && strings.Contains(qVal, "#lastresult") && len(results)-1 >= 0 {
						fmt.Printf("%s\n", results[len(results)-1][idx])
					}
				}

				color.HiBlue("%s", text)
			}
		}
		fmt.Println()
		results = append(results, result)
		if *isResultAsNode && len(selected) > 0 {
			currDoc = selected[0]
		}
	}
	fmt.Println()
}

func querySelect(queryXPATH bool, currDoc *html.Node, val string) (selected []*html.Node) {
	var err error
	if queryXPATH {
		htmlquery.Query(currDoc, val)
		selected, err = htmlquery.QueryAll(currDoc, val)
		if err != nil {
			panic(err)
		}
	} else {
		sel, err := cascadia.Parse(val)
		if err != nil {
			panic(err)
		}
		selected = cascadia.QueryAll(currDoc, sel)
	}
	return selected
}

// get the content of an html.Node
// https://stackoverflow.com/a/18275336
func collectText(n *html.Node, buf *bytes.Buffer, level int) {
	if n.Type == html.TextNode {
		if value, exists := headerMap[len(results)]; exists {
			if strs, ok := value.([]string); ok && level < len(strs) {
				buf.WriteString(strs[level] + ":" + n.Data + "\n")
			}
		} else {
			buf.WriteString(n.Data + "\n")
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf, level)
		//To correctly read the header names, increment the level by one after processing each child td element node.
		if c.Type == html.ElementNode && c.Data == "td" {
			level++
		}
	}
}
