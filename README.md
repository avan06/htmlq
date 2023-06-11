htmlq v1.0.0
====

`htmlq` is a `command-line` tool that allows you to query HTML using `CSS selectors` or `XPATH` and retrieve the corresponding `text content` (similar to JavaScript's `document.querySelector(query).textContent`).


Usage
====

```
usage: htmlq 1.0.0 [-h|--help] [-f|--file "<value>"] [-t|--text "<value>"]
                   [-u|--url "<value>"] [-x|--XPATH] [-a|--SelectorAll]
                   [-r|--ResultAsNode] [-l|--PrintLastResult] [-H|--headers
                   "<value>" [-H|--headers "<value>" ...]] [-v|--verbose]
                   -q|--querys "<value>" [-q|--querys "<value>" ...]

                   A command-line tool that allows you to query HTML using CSS
                   selectors or XPATH and retrieve the corresponding text
                   content (similar to JavaScript's
                   `document.querySelector(query).textContent`)

Arguments:

  -h  --help             Print help information
  -f  --file             Enter the relative or absolute path of the HTML file
  -t  --text             Enter the HTML text content
  -u  --url              Enter the URL of the HTML
  -x  --XPATH            Enable default XPATH query syntax instead of CSS
                         Selectors. Default: false
  -a  --SelectorAll      Enable the SelectorAll mechanism. When there are
                         multiple results for a single query, it will return
                         all the results (similar to querySelectorAll). By
                         default, it is not enabled, and only a single result
                         is returned (similar to querySelector).. Default:
                         false
  -r  --ResultAsNode     Enable using the Node from the previous query result
                         as the current query's root Node. Default: false
  -l  --PrintLastResult  Enable printing the content of the last result in the
                         output when using the "#lastresult" syntax in the
                         query. Default: false
  -H  --headers          You can input corresponding names for each result of a
                         single query using the format "#{serial
                         number}:header1Name;header2Name;header3Name;...". The
                         serial number represents the Nth query, starting from
                         zero.
  -v  --verbose          verbose
  -q  --querys           Enter a query value to retrieve the Node Text of the
                         target HTML. The default query method is CSS
                         Selectors, but it can be changed to XPATH using "-x"
                         or "--XPATH". Both query methods can be mixed. When a
                         query starts with "C:", it represents CSS Selectors,
                         and when it starts with "X:", it represents
                         XPATH.

	You can enter multiple query values separated by spaces and return all the
                         results.
	
	* When the "--ResultNode" is enabled, you can use the following special query
                         values to change the current position of the HTML
                         Node:
	  - Parent: Move the Node up one level
	  - NextSibling: Move the Node to the next sibling in the same level
	  - PrevSibling: Move the Node to the previous sibling in the same level
	  - FirstChild: Move the Node to the first child in the same level
	  - LastChild: Move the Node to the last child in the same level
	  - reset: Restore the Node to the root node of the original input
	
	* When using the "#lastresult" query syntax, the text content of the previous
                         query will automatically replace
                         "#lastresult".
	
	Example: --ResultNode -q Query1 Parent Query2 NextSibling Query3 LastChild
                         C:Query4 reset X:Query5
```

Examples
====

Note: In command-line, to break a line in Linux you use "\\", while in Windows you use "^".

```
htmlq --XPATH --SelectorAll -vv --PrintLastResult 
-t "<table class='ws-table-all notranslate'><tbody><tr><th>Tag</th><th>Description</th></tr><tr><td>&lt;table&gt;</td><td>Defines a table</td></tr><tr><td>&lt;th&gt;</td><td>Defines a header cell in a table</td></tr><tr><td>&lt;tr&gt;</td><td>Defines a row in a table</td></tr><tr><td>&lt;td&gt;</td><td>Defines a cell in a table</td></tr><tr><td>&lt;caption&gt;</td><td>Defines a table caption</td></tr><tr><td>&lt;colgroup&gt;</td><td>Specifies a group of one or more columns in a table for formatting</td></tr><tr><td>&lt;col&gt;</td><td>Specifies column properties for each column within a colgroup element</td></tr><tr><td>&lt;thead&gt;</td><td>Groups the header content in a table</td></tr><tr><td>&lt;tbody&gt;</td><td>Groups the body content in a table</td></tr><tr><td>&lt;tfoot&gt;</td><td>Groups the footer content in a table</td></tr></tbody></table>" 
--headers "#1:Tag1;Description1" --headers "#3:Tag3;Description3" 
-q "//table[contains(@class,'ws-table-all')]//tr[5]" 
"C:table.ws-table-all tr:nth-child(5)" 
"C:table.ws-table-all tr" 
"//table[contains(@class,'ws-table-all')]//tr[not(th)]" 
"//table[contains(@class,'ws-table-all')]//tr/td[1]" 
"//table[contains(@class,'ws-table-all')]//tr[contains(.,'#lastresult')]/td[2]"
```

```
0: //table[contains(@class,'ws-table-all')]//tr[5]
0-0:
<td>
Defines a cell in a table
```

```
1: table.ws-table-all tr:nth-child(5)
1-0:
Tag1:<td>
Description1:Defines a cell in a table
```

```
2: table.ws-table-all tr
2-0:
Tag
Description
2-1:
<table>
Defines a table
2-2:
<th>
Defines a header cell in a table
2-3:
<tr>
Defines a row in a table
2-4:
<td>
Defines a cell in a table
2-5:
<caption>
Defines a table caption
2-6:
<colgroup>
Specifies a group of one or more columns in a table for formatting
2-7:
<col>
Specifies column properties for each column within a colgroup element
2-8:
<thead>
Groups the header content in a table
2-9:
<tbody>
Groups the body content in a table
2-10:
<tfoot>
Groups the footer content in a table
```

```
3: //table[contains(@class,'ws-table-all')]//tr[not(th)]
3-0:
Tag3:<table>
Description3:Defines a table
3-1:
Tag3:<th>
Description3:Defines a header cell in a table
3-2:
Tag3:<tr>
Description3:Defines a row in a table
3-3:
Tag3:<td>
Description3:Defines a cell in a table
3-4:
Tag3:<caption>
Description3:Defines a table caption
3-5:
Tag3:<colgroup>
Description3:Specifies a group of one or more columns in a table for formatting
3-6:
Tag3:<col>
Description3:Specifies column properties for each column within a colgroup element
3-7:
Tag3:<thead>
Description3:Groups the header content in a table
3-8:
Tag3:<tbody>
Description3:Groups the body content in a table
3-9:
Tag3:<tfoot>
Description3:Groups the footer content in a table
```

```
4: //table[contains(@class,'ws-table-all')]//tr/td[1]
4-0:
<table>
4-1:
<th>
4-2:
<tr>
4-3:
<td>
4-4:
<caption>
4-5:
<colgroup>
4-6:
<col>
4-7:
<thead>
4-8:
<tbody>
4-9:
<tfoot>
```

```
5: //table[contains(@class,'ws-table-all')]//tr[contains(.,'#lastresult')]/td[2]
5-0:
<table>
Defines a table
5-1:
<th>
Defines a header cell in a table
5-2:
<tr>
Defines a row in a table
5-3:
<td>
Defines a cell in a table
5-4:
<caption>
Defines a table caption
5-5:
<colgroup>
Specifies a group of one or more columns in a table for formatting
5-6:
<col>
Specifies column properties for each column within a colgroup element
5-7:
<thead>
Groups the header content in a table
5-8:
<tbody>
Groups the body content in a table
5-9:
<tfoot>
Groups the footer content in a table
```


Credits
====

* [antchfx/htmlquery](https://github.com/antchfx/htmlquery)
* [ericchiang/css](https://github.com/ericchiang/css)
* [akamensky/argparse](https://github.com/akamensky/argparse)
* [fatih/color](https://github.com/fatih/color)

