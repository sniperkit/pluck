package cast

import (
	"math"

	"testing"
)

func TestFloat64FromFloat32(t *testing.T) {

	tests := []struct{
		Value float32
	}{
		{
			Value: float32(math.Inf(-1)),
		},
		{
			Value: -math.MaxFloat32,
		},
		{
			Value: -math.Pi,
		},
		{
			Value: -math.E,
		},
		{
			Value: -math.Sqrt2,
		},
		{
			Value: -1.0,
		},
		{
			Value: -math.Ln2,
		},
		{
			Value: -math.SmallestNonzeroFloat32,
		},
		{
			Value: 0.0,
		},
		{
			Value: math.SmallestNonzeroFloat32,
		},
		{
			Value: math.Ln2,
		},
		{
			Value: 1.0,
		},
		{
			Value: math.Sqrt2,
		},
		{
			Value: math.E,
		},
		{
			Value: math.Pi,
		},
		{
			Value: math.MaxFloat32,
		},
		{
			Value: float32(math.Inf(+1)),
		},



		{
			Value: float32(math.NaN()),
		},
	}

	{
		const numRand = 20
		for i:=0; i<numRand; i++ {
			test := struct{
				Value float32
			}{
				Value: randomness.Float32(),
			}
			tests = append(tests, test)

			test = struct{
				Value float32
			}{
				Value: -randomness.Float32(),
			}
			tests = append(tests, test)



			test = struct{
				Value float32
			}{
				Value: randomness.Float32() * math.MaxFloat32,
			}
			tests = append(tests, test)

			test = struct{
				Value float32
			}{
				Value: -randomness.Float32() * math.MaxFloat32,
			}
			tests = append(tests, test)



			test = struct{
				Value float32
			}{
				Value: randomness.Float32() * 999999999,
			}
			tests = append(tests, test)

			test = struct{
				Value float32
			}{
				Value: -randomness.Float32() * 999999999,
			}
			tests = append(tests, test)
		}
	}


	for testNumber, test := range tests {

		x, err := Float64(test.Value)
		if nil != err {
			t.Errorf("For test #%d, did not expect an error, but actually got one: (%T) %v", testNumber, err, err)
			continue
		}

		y := float32(x)

		if expected, actual := test.Value, y; expected != actual {
			if !(math.IsNaN(float64(expected)) && math.IsNaN(float64(actual))) {
				t.Errorf("For test #%d, expected %v, but actually got %v.", testNumber, expected, actual)
				continue
			}
		}
	}
}
