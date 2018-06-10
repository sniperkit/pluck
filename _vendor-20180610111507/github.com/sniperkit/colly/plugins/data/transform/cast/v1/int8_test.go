package cast

import (
	"math"

	"testing"
)

func TestInt8FromInt8(t *testing.T) {

	tests := []struct{
		Value int8
	}{
		{
			Value: math.MinInt8,
		},
		{
			Value: -1,
		},
		{
			Value: 0,
		},
		{
			Value: 1,
		},
		{
			Value: math.MaxInt8,
		},
	}

	{
		const numRand = 20
		for i:=0; i<numRand; i++ {
			test := struct{
				Value int8
			}{
				Value: int8(randomness.Int63n(math.MaxInt8)),
			}
			tests = append(tests, test)

			test = struct{
				Value int8
			}{
				Value: -int8(randomness.Int63n(-1*math.MinInt8)),
			}
			tests = append(tests, test)
		}
	}


	for testNumber, test := range tests {

		x, err := Int8(test.Value)
		if nil != err {
			t.Errorf("For test #%d, did not expect an error, but actually got one: (%T) %v", testNumber, err, err)
			continue
		}

		y := int8(x)

		if expected, actual := test.Value, y; expected != actual {
			t.Errorf("For test #%d, expected %v, but actually got %v.", testNumber, expected, actual)
			continue
		}
	}
}

func TestInt8FromInt8er(t *testing.T) {

	tests := []struct{
		Value    int8er
		Expected int8
	}{
		{
			Value: testInt8erMin(),
			Expected:   math.MinInt8,
		},
		{
			Value: testInt8erNegativeOne(),
			Expected:       -1,
		},
		{
			Value: testInt8erZero(),
			Expected:        0,
		},
		{
			Value: testInt8erOne(),
			Expected:        1,
		},
		{
			Value: testInt8erMax(),
			Expected:   math.MaxInt8,
		},
	}


	for testNumber, test := range tests {

		x, err := Int8(test.Value)
		if nil != err {
			t.Errorf("For test #%d, did not expect an error, but actually got one: (%T) %v", testNumber, err, err)
			continue
		}

		y := int8(x)

		if expected, actual := test.Expected, y; expected != actual {
			t.Errorf("For test #%d, expected %v, but actually got %v.", testNumber, expected, actual)
			continue
		}
	}
}
