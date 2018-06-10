package main

type SelectAxis string

const (
	AXIS_X SelectAxis = "y"
	AXIS_Y SelectAxis = "x"
)

type Selector struct {
	Enabled bool    `default:"true"` // Enabled select query
	Colums  AxisX   `json:"cols"`    // Columns to extract
	Rows    []AxisY `json:"rows"`    // Lines to extract
	isValid bool    // isValid syntax
	err     error   // occured errors
}

type AxisX struct {
	Enabled  bool `default:"true"`
	List     []Column
	isRange  bool
	isList   bool
	isUnique bool
	err      error // occured errors
}

type Column struct {
	Enabled bool `default:"true"`
	Name    string
	Rename  string
	isIndex bool
	err     error // occured errors
}

type AxisY struct {
	Enabled bool `default:"true"`
	Index   int
	err     error // occured errors
}
