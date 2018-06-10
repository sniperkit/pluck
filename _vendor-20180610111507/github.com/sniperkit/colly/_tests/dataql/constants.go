package main

const (
	RAW_TEXT = ` package main

// On every a element which has tabular format data call callback
// Notes:
c.OnTAB("0:0", func(e *colly.TABElement) {

	// col[:], rows[:]
	OnTAB("col[:], rows[:]", func(e *colly.TABElement) {})

// col[2], rows[1,6,8]
OnTAB("col[2], rows[1,6,8]", func(e *colly.TABElement) {})

// cols["name", "fullname"],rows[:10]
OnTAB("cols["name", "fullname"],rows[:10]", func(e *colly.TABElement) {})

// col["name", "fullname"], row[4]
OnTAB("cols["name", "fullname"],row[4]", func(e *colly.TABElement) {})

// cols[:], rows[1:7]
OnTAB("cols[:],rows[1:7]", func(e *colly.TABElement) {})

		// cols[1:], rows[:10]
		OnTAB("cols[1:],rows[:10]", func(e *colly.TABElement) {})

// cols[1:7], rows[1:]
OnTAB("cols[1:7],rows[1:]", func(e *colly.TABElement) {})

		// cols[:10], row[1,5,7]
		OnTAB("cols[:10],row[1,5,7]", func(e *colly.TABElement) {})

// row[1]
OnTAB("row[1]", func(e *colly.TABElement) {})

		// rows[1,2]
		OnTAB("row[1,2]", func(e *colly.TABElement) {})

// rows[1,5,7]
OnTAB("row[1,5,7]", func(e *colly.TABElement) {})

		// rows[:]
		OnTAB("rows[:]", func(e *colly.TABElement) {})

// rows[1:]
OnTAB("rows[1:]", func(e *colly.TABElement) {})

		// rows[1:7]
		OnTAB("rows[1:7]", func(e *colly.TABElement) {})

// rows[:10]
OnTAB("rows[:10]", func(e *colly.TABElement) {})

// col[2]
OnTAB("col[2]", func(e *colly.TABElement) {})

	// cols["name", "fullname"]
OnTAB("cols["name", "fullname"]", func(e *colly.TABElement) {})

// col["name", "fullname"]
OnTAB("col["name", "fullname"]", func(e *colly.TABElement) {})

						// cols[:]
						OnTAB("cols[:]", func(e *colly.TABElement) {})

						// cols[1:]
						OnTAB("cols[1:]", func(e *colly.TABElement) {})

						// cols[1:7]
						OnTAB("cols[1:7]", func(e *colly.TABElement) {})

// cols[:10]
OnTAB("cols[:10]", func(e *colly.TABElement) {})
`
)
