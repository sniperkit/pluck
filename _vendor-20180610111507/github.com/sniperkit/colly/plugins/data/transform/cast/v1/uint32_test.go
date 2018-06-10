package cast

import (
	"math"

	"testing"
)

func TestUint32FromUint32(t *testing.T) {

	tests := []struct{
		Value uint32
	}{
		{
			Value: 0,
		},
		{
			Value: 1,
		},
		{
			Value: math.MaxUint32,
		},
	}

	{
		const numRand = 20
		for i:=0; i<numRand; i++ {
			test := struct{
				Value uint32
			}{
				Value: uint32(randomness.Int63n(math.MaxUint32)),
			}
			tests = append(tests, test)
		}
	}


	for testNumber, test := range tests {

		x, err := Uint32(test.Value)
		if nil != err {
			t.Errorf("For test #%d, did not expect an error, but actually got one: (%T) %v", testNumber, err, err)
			continue
		}

		y := uint32(x)

		if expected, actual := test.Value, y; expected != actual {
			t.Errorf("For test #%d, expected %v, but actually got %v.", testNumber, expected, actual)
			continue
		}
	}
}

func TestUint32FromUint32er(t *testing.T) {

	tests := []struct{
		Value    uint32er
		Expected uint32
	}{
		{
			Value: testUint32erZero(),
			Expected:          0,
		},
		{
			Value: testUint32erOne(),
			Expected:          1,
		},
		{
			Value: testUint32erMax(),
			Expected:     math.MaxUint32,
		},
	}


	for testNumber, test := range tests {

		x, err := Uint32(test.Value)
		if nil != err {
			t.Errorf("For test #%d, did not expect an error, but actually got one: (%T) %v", testNumber, err, err)
			continue
		}

		y := uint32(x)

		if expected, actual := test.Expected, y; expected != actual {
			t.Errorf("For test #%d, expected %v, but actually got %v.", testNumber, expected, actual)
			continue
		}
	}
}

func TestUint32FromUint16(t *testing.T) {

	tests := []struct{
		Value uint16
	}{
		{
			Value: 0,
		},
		{
			Value: 1,
		},
		{
			Value: math.MaxUint16,
		},
	}

	{
		const numRand = 20
		for i:=0; i<numRand; i++ {
			test := struct{
				Value uint16
			}{
				Value: uint16(randomness.Int63n(math.MaxUint16)),
			}
			tests = append(tests, test)
		}
	}


	for testNumber, test := range tests {

		x, err := Uint32(test.Value)
		if nil != err {
			t.Errorf("For test #%d, did not expect an error, but actually got one: (%T) %v", testNumber, err, err)
			continue
		}

		y := uint16(x)

		if expected, actual := test.Value, y; expected != actual {
			t.Errorf("For test #%d, expected %v, but actually got %v.", testNumber, expected, actual)
			continue
		}
	}
}

func TestUint32FromUint16er(t *testing.T) {

	tests := []struct{
		Value    uint16er
		Expected uint16
	}{
		{
			Value: testUint16erZero(),
			Expected:         0,
		},
		{
			Value: testUint16erOne(),
			Expected:         1,
		},
		{
			Value: testUint16erMax(),
			Expected:    math.MaxUint16,
		},
	}


	for testNumber, test := range tests {

		x, err := Uint32(test.Value)
		if nil != err {
			t.Errorf("For test #%d, did not expect an error, but actually got one: (%T) %v", testNumber, err, err)
			continue
		}

		y := uint16(x)

		if expected, actual := test.Expected, y; expected != actual {
			t.Errorf("For test #%d, expected %v, but actually got %v.", testNumber, expected, actual)
			continue
		}
	}
}

func TestUint32FromUint8(t *testing.T) {

	tests := []struct{
		Value uint8
	}{
		{
			Value: 0,
		},
		{
			Value: 1,
		},
		{
			Value: math.MaxUint8,
		},
	}

	{
		const numRand = 20
		for i:=0; i<numRand; i++ {
			test := struct{
				Value uint8
			}{
				Value: uint8(randomness.Int63n(math.MaxUint8)),
			}
			tests = append(tests, test)
		}
	}


	for testNumber, test := range tests {

		x, err := Uint32(test.Value)
		if nil != err {
			t.Errorf("For test #%d, did not expect an error, but actually got one: (%T) %v", testNumber, err, err)
			continue
		}

		y := uint8(x)

		if expected, actual := test.Value, y; expected != actual {
			t.Errorf("For test #%d, expected %v, but actually got %v.", testNumber, expected, actual)
			continue
		}
	}
}

func TestUint32FromUint8er(t *testing.T) {

	tests := []struct{
		Value    uint8er
		Expected uint8
	}{
		{
			Value: testUint8erZero(),
			Expected:         0,
		},
		{
			Value: testUint8erOne(),
			Expected:         1,
		},
		{
			Value: testUint8erMax(),
			Expected:    math.MaxUint8,
		},
	}


	for testNumber, test := range tests {

		x, err := Uint32(test.Value)
		if nil != err {
			t.Errorf("For test #%d, did not expect an error, but actually got one: (%T) %v", testNumber, err, err)
			continue
		}

		y := uint8(x)

		if expected, actual := test.Expected, y; expected != actual {
			t.Errorf("For test #%d, expected %v, but actually got %v.", testNumber, expected, actual)
			continue
		}
	}
}
