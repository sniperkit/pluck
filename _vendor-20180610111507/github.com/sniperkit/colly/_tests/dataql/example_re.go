package main

import (
	"fmt"

	regex "github.com/dakerfp/re"
)

func example_re() {
	reTest := regex.Regex(
		regex.Group("dividend", regex.Digits),
		regex.Then("/"),
		regex.Group("divisor", regex.Digits),
	)

	m := reTest.FindStringSubmatch("4/3")
	fmt.Println(m[1]) // > 4
	fmt.Println(m[2]) // > 3
}
