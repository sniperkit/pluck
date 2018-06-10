package main

import (
	"fmt"
	// "os"

	exprel "layeh.com/exprel"
)

func parse_example_exprel(queries ...string) {
	for _, query := range queries {
		expr, err := exprel.Parse(query)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}
		result, err := expr.Evaluate(exprel.Base)
		if err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Println("result", result)
		}
	}
}
