# Colly - Export format-agnostic tabular dataset

## Usage

Creates a dataset and populate it:

```go
	c.OnTAB("0:0", func(e *colly.TABElement) {


		ds := e.NewDataset([]string{"firstName", "lastName"})
	}


```

Add new rows:
```go
ds.Append([]interface{}{"John", "Adams"})
ds.AppendValues("George", "Washington")
```

Add new columns:
```go
ds.AppendColumn("age", []interface{}{90, 67})
ds.AppendColumnValues("sex", "male", "male")
```

Add a dynamic column, by passing a function which has access to the current row, and must
return a value:
```go
func lastNameLen(row []interface{}) interface{} {
	return len(row[1].(string))
}
ds.AppendDynamicColumn("lastName length", lastNameLen)
ds.CSV()
// >>
// firstName, lastName, age, sex, lastName length
// John, Adams, 90, male, 5
// George, Washington, 67, male, 10
```

Delete rows:
```go
ds.DeleteRow(1) // starts at 0
```

Delete columns:
```go
ds.DeleteColumn("sex")
```

Get a row or multiple rows:
```go
row, _ := ds.Row(0)
fmt.Println(row["firstName"]) // George

rows, _ := ds.Rows(0, 1)
fmt.Println(rows[0]["firstName"]) // George
fmt.Println(rows[1]["firstName"]) // Thomas
```

Slice a Dataset:
```go
newDs, _ := ds.Slice(1, 5) // returns a fresh Dataset with rows [1..5[
```


## Filtering

You can add **tags** to rows by using a specific `Dataset` method. This allows you to filter your `Dataset` later. This can be useful to separate rows of data based on arbitrary criteria (e.g. origin) that you don’t want to include in your `Dataset`.
```go
ds := NewDataset([]string{"Maker", "Model"})
ds.AppendTagged([]interface{}{"Porsche", "911"}, "fast", "luxury")
ds.AppendTagged([]interface{}{"Skoda", "Octavia"}, "family")
ds.AppendTagged([]interface{}{"Ferrari", "458"}, "fast", "luxury")
ds.AppendValues("Citroen", "Picasso")
ds.AppendValues("Bentley", "Continental")
ds.Tag(4, "luxury") // Bentley
ds.AppendValuesTagged("Aston Martin", "DB9", /* these are tags */ "fast", "luxury")
```

Filtering the `Dataset` is possible by calling `Filter(column)`:
```go
luxuryCars, err := ds.Filter("luxury").CSV()
fmt.Println(luxuryCars)
// >>>
// Maker,Model
// Porsche,911
// Ferrari,458
// Bentley,Continental
// Aston Martin,DB9
```

```go
fastCars, err := ds.Filter("fast").CSV()
fmt.Println(fastCars)
// >>>
// Maker,Model
// Porsche,911
// Ferrari,458
// Aston Martin,DB9
```

Tags at a specific row can be retrieved by calling `Dataset.Tags(index int)`

## Sorting

Datasets can be sorted by a specific column.
```go
ds := NewDataset([]string{"Maker", "Model", "Year"})
ds.AppendValues("Porsche", "991", 2012)
ds.AppendValues("Skoda", "Octavia", 2011)
ds.AppendValues("Ferrari", "458", 2009)
ds.AppendValues("Citroen", "Picasso II", 2013)
ds.AppendValues("Bentley", "Continental GT", 2003)

sorted, err := ds.Sort("Year").CSV()
fmt.Println(sorted)
// >>
// Maker, Model, Year
// Bentley, Continental GT, 2003
// Ferrari, 458, 2009
// Skoda, Octavia, 2011
// Porsche, 991, 2012
// Citroen, Picasso II, 2013
```

## Constraining

Datasets can have columns constrained by functions and further checked if valid.
```go
ds := NewDataset([]string{"Maker", "Model", "Year"})
ds.AppendValues("Porsche", "991", 2012)
ds.AppendValues("Skoda", "Octavia", 2011)
ds.AppendValues("Ferrari", "458", 2009)
ds.AppendValues("Citroen", "Picasso II", 2013)
ds.AppendValues("Bentley", "Continental GT", 2003)

ds.ConstrainColumn("Year", func(val interface{}) bool { return val.(int) > 2008 })
ds.ValidFailFast() // false
if !ds.Valid() { // validate the whole dataset, errors are retrieved in Dataset.ValidationErrors
	ds.ValidationErrors[0] // Row: 4, Column: 2
}
```

A Dataset with constrained columns can be filtered to keep only the rows satisfying the constraints.
```go
valid := ds.ValidSubset().Tabular("simple") // Cars after 2008
fmt.Println(valid)
```

Will output:
```
------------  ---------------  ---------
      Maker            Model       Year
------------  ---------------  ---------
    Porsche              991       2012

      Skoda          Octavia       2011

    Ferrari              458       2009

    Citroen       Picasso II       2013
------------  ---------------  ---------
```

```go
invalid := ds.InvalidSubset().Tabular("simple") // Cars before 2008
fmt.Println(invalid)
```

Will output:
```
------------  -------------------  ---------
      Maker                Model       Year
------------  -------------------  ---------
    Bentley       Continental GT       2003
------------  -------------------  ---------
```

## Loading

### JSON
```go
ds, _ := LoadJSON([]byte(`[
  {"age":90,"firstName":"John","lastName":"Adams"},
  {"age":67,"firstName":"George","lastName":"Washington"},
  {"age":83,"firstName":"Henry","lastName":"Ford"}
]`))
```
