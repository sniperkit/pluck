# Colly - Export format-agnostic tabular dataset

## Exports

### Exportable

Any of the following export format returns an `*Exportable` which means you can use:
- `Bytes()` to get the content as a byte array
- `String()` to get the content as a string
- `WriteTo(io.Writer)` to write the content to an `io.Writer`
- `WriteFile(filename string, perm os.FileMode)` to write to a file

It avoids unnecessary conversion between `string` and `[]byte` to output/write/whatever.
Thanks to [@figlief](https://github.com/figlief) for the proposition. 

### JSON
```go
json, _ := ds.JSON()
fmt.Println(json)
```

Will output:
```json
[{"age":90,"firstName":"John","lastName":"Adams"},{"age":67,"firstName":"George","lastName":"Washington"},{"age":83,"firstName":"Henry","lastName":"Ford"}]
```

### XML
```go
xml, _ := ds.XML()
fmt.Println(xml)
```

Will ouput:
```xml
<dataset>
 <row>
   <age>90</age>
   <firstName>John</firstName>
   <lastName>Adams</lastName>
 </row>  <row>
   <age>67</age>
   <firstName>George</firstName>
   <lastName>Washington</lastName>
 </row>  <row>
   <age>83</age>
   <firstName>Henry</firstName>
   <lastName>Ford</lastName>
 </row>
</dataset>
```

### CSV
```go
csv, _ := ds.CSV()
fmt.Println(csv)
```

Will ouput:
```csv
firstName,lastName,age
John,Adams,90
George,Washington,67
Henry,Ford,83
```

### TSV
```go
tsv, _ := ds.TSV()
fmt.Println(tsv)
```

Will ouput:
```tsv
firstName lastName  age
John  Adams  90
George  Washington  67
Henry Ford 83
```

### YAML
```go
yaml, _ := ds.YAML()
fmt.Println(yaml)
```

Will ouput:
```yaml
- age: 90
  firstName: John
  lastName: Adams
- age: 67
  firstName: George
  lastName: Washington
- age: 83
  firstName: Henry
  lastName: Ford
```

### HTML
```go
html := ds.HTML()
fmt.Println(html)
```

Will output:
```html
<table class="table table-striped">
	<thead>
		<tr>
			<th>firstName</th>
			<th>lastName</th>
			<th>age</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td>George</td>
			<td>Washington</td>
			<td>90</td>
		</tr>
		<tr>
			<td>Henry</td>
			<td>Ford</td>
			<td>67</td>
		</tr>
		<tr>
			<td>Foo</td>
			<td>Bar</td>
			<td>83</td>
		</tr>
	</tbody>
</table>
```

### XLSX
```go
xlsx, _ := ds.XLSX()
fmt.Println(xlsx)
// >>>
// binary content
xlsx.WriteTo(...)
```

### ASCII

#### Grid format
```go
ascii := ds.Tabular("grid" /* tablib.TabularGrid */)
fmt.Println(ascii)
```

Will output:
```
+--------------+---------------+--------+
|    firstName |      lastName |    age |
+==============+===============+========+
|       George |    Washington |     90 |
+--------------+---------------+--------+
|        Henry |          Ford |     67 |
+--------------+---------------+--------+
|          Foo |           Bar |     83 |
+--------------+---------------+--------+
```

#### Simple format
```go
ascii := ds.Tabular("simple" /* tablib.TabularSimple */)
fmt.Println(ascii)
```

Will output:
```
--------------  ---------------  --------
    firstName         lastName       age
--------------  ---------------  --------
       George       Washington        90

        Henry             Ford        67

          Foo              Bar        83
--------------  ---------------  --------
```

#### Condensed format
```go
ascii := ds.Tabular("condensed" /* tablib.TabularCondensed */)
fmt.Println(ascii)
```

Similar to simple but with less line feed:
```
--------------  ---------------  --------
    firstName         lastName       age
--------------  ---------------  --------
       George       Washington        90
        Henry             Ford        67
          Foo              Bar        83
--------------  ---------------  --------
```

### Markdown

Markdown tables are similar to the Tabular condensed format, except that they have
pipe characters separating columns.

```go
mkd := ds.Markdown() // or
mkd := ds.Tabular("markdown" /* tablib.TabularMarkdown */)
fmt.Println(mkd)
```

Will output:
```
|     firstName   |       lastName    |    gpa  |
| --------------  | ---------------   | ------- |
|          John   |          Adams    |     90  |
|        George   |     Washington    |     67  |
|        Thomas   |      Jefferson    |     50  |
```

Which equals to the following when rendered as HTML:

|     firstName   |       lastName    |    gpa  |
| --------------  | ---------------   | ------- |
|          John   |          Adams    |     90  |
|        George   |     Washington    |     67  |
|        Thomas   |      Jefferson    |     50  |

### MySQL
```go
sql := ds.MySQL()
fmt.Println(sql)
```

Will output:
```sql
CREATE TABLE IF NOT EXISTS presidents
(
	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	firstName VARCHAR(9),
	lastName VARCHAR(8),
	gpa DOUBLE
);

INSERT INTO presidents VALUES(1, 'Jacques', 'Chirac', 88);
INSERT INTO presidents VALUES(2, 'Nicolas', 'Sarkozy', 98);
INSERT INTO presidents VALUES(3, 'François', 'Hollande', 34);

COMMIT;
```

Numeric (`uint`, `int`, `float`, ...) are stored as `DOUBLE`, `string`s as `VARCHAR` with width set to the length of the longest string in the column, and `time.Time`s are stored as `TIMESTAMP`.

### Postgres
```go
sql := ds.Postgres()
fmt.Println(sql)
```

Will output:
```sql
CREATE TABLE IF NOT EXISTS presidents
(
	id SERIAL PRIMARY KEY,
	firstName TEXT,
	lastName TEXT,
	gpa NUMERIC
);

INSERT INTO presidents VALUES(1, 'Jacques', 'Chirac', 88);
INSERT INTO presidents VALUES(2, 'Nicolas', 'Sarkozy', 98);
INSERT INTO presidents VALUES(3, 'François', 'Hollande', 34);

COMMIT;
```

Numerics (`uint`, `int`, `float`, ...) are stored as `NUMERIC`, `string`s as `TEXT` and `time.Time`s are stored as `TIMESTAMP`.

## Databooks

This is an example of how to use Databooks.

```go
db := NewDatabook()
// or loading a JSON content
db, err := LoadDatabookJSON([]byte(`...`))
// or a YAML content
db, err := LoadDatabookYAML([]byte(`...`))

// a dataset of presidents
presidents, _ := LoadJSON([]byte(`[
  {"Age":90,"First name":"John","Last name":"Adams"},
  {"Age":67,"First name":"George","Last name":"Washington"},
  {"Age":83,"First name":"Henry","Last name":"Ford"}
]`))

// a dataset of cars
cars := NewDataset([]string{"Maker", "Model", "Year"})
cars.AppendValues("Porsche", "991", 2012)
cars.AppendValues("Skoda", "Octavia", 2011)
cars.AppendValues("Ferrari", "458", 2009)
cars.AppendValues("Citroen", "Picasso II", 2013)
cars.AppendValues("Bentley", "Continental GT", 2003)

// add the sheets to the Databook
db.AddSheet("Cars", cars.Sort("Year"))
db.AddSheet("Presidents", presidents.SortReverse("Age"))

fmt.Println(db.JSON())
```

Will output the following JSON representation of the Databook:
```json
[
  {
    "title": "Cars",
    "data": [
      {"Maker":"Bentley","Model":"Continental GT","Year":2003},
      {"Maker":"Ferrari","Model":"458","Year":2009},
      {"Maker":"Skoda","Model":"Octavia","Year":2011},
      {"Maker":"Porsche","Model":"991","Year":2012},
      {"Maker":"Citroen","Model":"Picasso II","Year":2013}
    ]
  },
  {
    "title": "Presidents",
    "data": [
      {"Age":90,"First name":"John","Last name":"Adams"},
      {"Age":83,"First name":"Henry","Last name":"Ford"},
      {"Age":67,"First name":"George","Last name":"Washington"}
    ]
  }
]
```