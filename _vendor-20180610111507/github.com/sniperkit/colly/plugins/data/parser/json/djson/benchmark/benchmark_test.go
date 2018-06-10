package benchmark

import (
	"encoding/json"
	"testing"

	"github.com/Jeffail/gabs"
	"github.com/a8m/djson"
	"github.com/antonholmquist/jason"
	"github.com/bitly/go-simplejson"
	"github.com/mreiferson/go-ujson"
	"github.com/ugorji/go/codec"
)

func BenchmarkEncodingJsonParser(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := make(map[string]interface{})
			json.Unmarshal(smallFixture, &data)
		}
	})

	b.Run("medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := make(map[string]interface{})
			json.Unmarshal(mediumFixture, &data)
		}
	})

	b.Run("large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := make(map[string]interface{})
			json.Unmarshal(largeFixture, &data)
		}
	})
}

func BenchmarkUgorjiParser(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			decoder := codec.NewDecoderBytes(smallFixture, new(codec.JsonHandle))
			var v interface{}
			decoder.Decode(&v)
		}
	})

	b.Run("medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			decoder := codec.NewDecoderBytes(mediumFixture, new(codec.JsonHandle))
			var v interface{}
			decoder.Decode(&v)
		}
	})

	b.Run("large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			decoder := codec.NewDecoderBytes(largeFixture, new(codec.JsonHandle))
			var v interface{}
			decoder.Decode(&v)
		}
	})
}

func BenchmarkJasonParser(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			jason.NewObjectFromBytes(smallFixture)
		}
	})

	b.Run("medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			jason.NewObjectFromBytes(mediumFixture)
		}
	})

	b.Run("large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			jason.NewObjectFromBytes(largeFixture)
		}
	})
}

func BenchmarkSimpleJsonParser(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			simplejson.NewJson(smallFixture)
		}
	})

	b.Run("medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			simplejson.NewJson(mediumFixture)
		}
	})

	b.Run("large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			simplejson.NewJson(largeFixture)
		}
	})
}

func BenchmarkGabsParser(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gabs.ParseJSON(smallFixture)
		}
	})

	b.Run("medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gabs.ParseJSON(mediumFixture)
		}
	})

	b.Run("large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gabs.ParseJSON(largeFixture)
		}
	})
}

func BenchmarkUJsonParser(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ujson.NewFromBytes(smallFixture)
		}
	})

	b.Run("medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ujson.NewFromBytes(mediumFixture)
		}
	})

	b.Run("large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ujson.NewFromBytes(largeFixture)
		}
	})
}

func BenchmarkDJsonParser(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			djson.DecodeObject(smallFixture)
		}
	})

	b.Run("medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			djson.DecodeObject(mediumFixture)
		}
	})

	b.Run("large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			djson.DecodeObject(largeFixture)
		}
	})
}

func BenchmarkDJsonAllocString(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dec := djson.NewDecoder(smallFixture)
			dec.AllocString()
			dec.DecodeObject()
		}
	})

	b.Run("medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dec := djson.NewDecoder(mediumFixture)
			dec.AllocString()
			dec.DecodeObject()
		}
	})

	b.Run("large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dec := djson.NewDecoder(largeFixture)
			dec.AllocString()
			dec.DecodeObject()
		}
	})
}

/*
// This is not part of the benchmark test cases;
// Trying to show the preformence when translating the jsonparser's
// result into map[string]interface{}
// import "github.com/buger/jsonparser"
func BenchmarkJsonparserParser(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m := make(map[string]interface{})
			jsonparser.ObjectEach(smallFixture, func(k, v []byte, vt jsonparser.ValueType, o int) error {
				if vt == jsonparser.Number {
					m[string(k)], _ = strconv.ParseFloat(string(v), 64)
				} else {
					m[string(k)] = string(v)
				}
				return nil
			})
		}
	})
}
*/
