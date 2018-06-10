package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	// helpers
	pp "github.com/sniperkit/colly/plugins/app/debug/pp"
	iter "github.com/sniperkit/colly/plugins/data/structure/int/iter"
)

var SELECTOR_QUERY_REGEX *regexp.Regexp

func parse_selector_queries(queries []string) {
	SELECTOR_QUERY_REGEX = regexp.MustCompile(SELECT_ALPHA_NUMERIC)
	for _, query := range queries {
		parse_selector_query(query)
	}
}

func stringsSliceToIntsSlice(sliceStr []string) (sliceInt []int, errs []string) {
	for _, keyStr := range sliceStr {
		if keyInt, err := strconv.Atoi(keyStr); err == nil {
			sliceInt = append(sliceInt, keyInt)
		} else {
			errs = append(errs, err.Error())
		}
	}
	return
}

func range2List(lower, upper, limit int) (list []int) {
	if limit == 0 {
		limit = upper
	}
	if upper > limit {
		upper = limit
	}
	fmt.Println("lower=", lower, "upper=", upper, "limit=", limit)
	for i := lower; i <= upper; i++ {
		list = append(list, i)
	}
	return list
}

func iterRange2List(lower, upper, limit int) (list []int) {
	if limit == 0 {
		limit = upper
	}
	if upper > limit {
		upper = limit
	}
	fmt.Println("lower=", lower, "upper=", upper, "limit=", limit)
	for i := range iter.N(upper) {
		list = append(list, i)
	}
	return list
}

func parse_selector_query(queryStr string) {

	// remove extra spaces in query string
	queryStr = strings.Replace(queryStr, " ", "", -1)
	queryStr = strings.Replace(queryStr, "\"", "", -1)

	// sanitize ?!
	// queryStr = sanitize_selector_query(queryStr)

	selectParts := SELECTOR_QUERY_REGEX.FindAllStringSubmatch(queryStr, SELECTOR_SIMPLE_ARGS_MAX)
	fmt.Println("\n------------------------------------------------------------------------------------------------------------------------")
	// fmt.Printf("extract for `%s\n", queryStr)

	// Check if we have at least one selection axis if wildcard is not defined
	if len(selectParts) < SELECTOR_SIMPLE_ARGS_MIN && queryStr != SELECTOR_SIMPLE_WILDCARD {
		fmt.Println("An error occured while parsing the selection query syntax.")
	}

	// Loop over matched selector query
	for _, selectPart := range selectParts {

		// define error flag
		var isError bool

		// define selection axis
		var selectAxis string
		switch selectPart[2] {
		case "row", "rows":
			selectAxis = "y"
		case "col", "cols":
			selectAxis = "x"
		default:
			isError = true
		}

		if isError {
			fmt.Println("error: unable to detect the select query axis... selectPart=", selectPart)
			continue
		}

		// get matched patterns parts
		selectList := strings.Split(selectPart[3], SELECTOR_SIMPLE_LIST_SEP)
		selectRange := strings.Split(selectPart[3], SELECTOR_SIMPLE_RANGE_SEP)

		var selectIndices, selectIndicesIter []int // declare vars for by indices
		var selectNames []string                   // declare vars for by names

		// declare selection type flags
		var isList, isRange, isUnique, isUniqueInt, isUniqueStr bool

		// declare lower, upper and cap value for ranges
		var rlower, rupper, rcap, uidx int

		// declare uniqueKey
		var ukey string

		// check matched patterns results
		switch {

		// select unique key can be either only numeric or only alphanumeric...
		case selectRange[0] == selectList[0]:

			// isnumeric
			switch {
			case IsNumeric(selectRange[0]):
				// isNumeric = true
				isUniqueInt = true
				if key, err := strconv.Atoi(selectList[0]); err == nil {
					uidx = key
					isUnique = true
				} else {
					isError = true
				}

			case IsAlphaNumeric(selectRange[0]):
				// isAlphanum = true
				isUniqueStr = true
				ukey = selectRange[0]

			default:
				isError = true
			}
			pp.Printf("select unique key: isUniqueInt='%s', isUniqueStr='%s' \n", isUniqueInt, isUniqueStr)

			// s := strconv.Itoa(-42)

		// Important! select list can be either only numeric or only alphanumeric...
		case len(selectList) >= 2 && len(selectRange) == 1:
			// check if numeric
			// selectIndices = strings.Split(selectList, ",")

			// listRaw = selectList
			selectNames = selectList

			// check if alphanumeric
			// selectNames = strings.Split(selectList, ",")
			var errs []string
			selectIndices, errs = stringsSliceToIntsSlice(selectList)

			switch {
			case len(errs) == len(selectList):
				selectIndices = []int{}
				// fmt.Println("detected list of alphanumeric keys")

			case len(errs) > 0 && len(errs) < len(selectList):
				fmt.Println("error occured while converting []string to []int.", selectList)

			case len(errs) == 0:
				selectNames = []string{}

			}

			// Get the list slices
			isList = true

		// Important! select ranges are only numeric...
		case (len(selectRange) >= 2 && len(selectRange) < 3) && len(selectList) == 1:
			switch {
			case selectRange[0] == "":
				rlower = 0
			default:
				if keyInt, err := strconv.Atoi(selectRange[0]); err == nil {
					rlower = keyInt
				} else {
					isError = true
				}
			}

			switch {
			case selectRange[1] == "":
				rupper = 0 // dataset length...
			default:
				if keyInt, err := strconv.Atoi(selectRange[1]); err == nil {
					rupper = keyInt
				} else {
					isError = true
				}
			}

			if rupper < rlower {
				isError = true
			}

			if len(selectRange) == 3 {
				switch {
				case selectRange[2] == "":
					if rupper > 0 {
						rcap = rupper // dataset length...
					}

				default:
					if keyInt, err := strconv.Atoi(selectRange[2]); err == nil {
						rcap = keyInt
					} else {
						isError = true
					}
				}
			}

			if rcap == 0 && rupper > 0 && rlower <= rupper {
				rcap = rupper
			}

			if rlower <= rupper && rupper > 0 {
				fmt.Println("attempt to generate the index list between the range...")
				selectIndices = range2List(rlower, rupper, rcap)
				// selectIndicesIter = iterRange2List(rlower, rupper, rcap)
			}

			isRange = true

		case len(selectList) >= 2 && len(selectRange) > 1:
			fallthrough

		default:
			isError = true
		}

		pp.Printf("queryStrs=\"%s\" \n", queryStr)
		pp.Printf("selectAxis=\"%s\", selectStr=`%s`, isList: \"%s\", isRange: \"%s\", isUnique=\"%s\", isError: \"%s\" \n", selectAxis, selectPart[3], isList, isRange, isUnique, isError)

		switch {
		case isRange:
			pp.Println("lowerRange=", rlower, "upperRange=", rupper, "capRange=", rcap, "selectRange.Length=", len(selectRange))
			pp.Println("len(selectIndices)=", len(selectIndices), " len(selectIndicesIter)", len(selectIndicesIter))
			if len(selectIndices) > 0 {
				pp.Println("selectIndices=", selectIndices)
			}
			if len(selectIndicesIter) > 0 {
				pp.Println("selectIndicesIter=", selectIndicesIter)
			}

		case isList:
			pp.Println("selectList.Length=", len(selectList))
			if len(selectIndices) > 0 {
				pp.Println("selectIndices=", selectIndices)
			}
			if len(selectNames) > 0 {
				pp.Println("selectNames=", selectNames)
			}
		case isUnique:
			pp.Printf("isUniqueStr=\"%s\", isUniqueInt=\"%s\" , uidx=\"%s\", ukey=\"%s\"", isUniqueStr, isUniqueInt, uidx, ukey)

		default:
			isError = true
		}

		isDebug := false
		if queryStr == "cols[1:7]" && isDebug {
			os.Exit(1)
		}

		if isError {
			continue
		}

		// check if unique key
		// check if keys are numeric only
		// check if keys are alphanumeric

		if isDebug {
			pp.Println("selectIndices=", selectIndices)
			pp.Println("selectNames=", selectNames)
		}

	}

	fmt.Println("")

}
