package tablib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	// plugins - json parsers
	/*
		gjson "github.com/sniperkit/colly/plugins/data/parser/json/gson"
		djson "github.com/sniperkit/colly/plugins/data/parser/json/djson"
		easyjson "github.com/sniperkit/colly/plugins/data/parser/json/easyjson"
		fastjson "github.com/sniperkit/colly/plugins/data/parser/json/fastjson"
		ffjson "github.com/sniperkit/colly/plugins/data/parser/json/ffjson"
		gabs "github.com/sniperkit/colly/plugins/data/parser/json/gabs"
		jsnm "github.com/sniperkit/colly/plugins/data/parser/json/jsnm"
		jsonez "github.com/sniperkit/colly/plugins/data/parser/json/jsonez"
		jsonparser "github.com/sniperkit/colly/plugins/data/parser/json/jsonparser"
		jsonstream "github.com/sniperkit/colly/plugins/data/parser/json/jsonstream"
		lazyjson "github.com/sniperkit/colly/plugins/data/parser/json/lazyjson"
	*/
	// gjson "github.com/sniperkit/colly/plugins/data/parser/json/gjson"    // fork
	mxj "github.com/sniperkit/colly/plugins/data/parser/json/mxj/v2/pkg" // fork v2
	// mxj "github.com/sniperkit/colly/plugins/data/transform/mxj/master" // latest commit
	// mxj "github.com/sniperkit/colly/plugins/data/transform/mxj/v1" // fork v1
	// dev helpers
	// pp "github.com/sniperkit/colly/plugins/app/debug/pp"
)

// LoadJSON loads a dataset from a JSON source.
// - Default pkg "encoding/json"
func LoadJSON(jsonContent []byte) (*Dataset, error) {

	var input []map[string]interface{}
	d := json.NewDecoder(strings.NewReader(string(jsonContent)))
	d.UseNumber()
	if err := d.Decode(&input); err != nil {
		return nil, err
	}

	// if err := json.Unmarshal(jsonContent, &input); err != nil {
	// 	return nil, err
	// }

	return internalLoadFromDict(input)
}

func convertInterface(input map[string]interface{}) []interface{} {
	results := make([]interface{}, len(input))
	for _, result := range input {
		resultSlice := result.(interface{})
		results = append(results, resultSlice)
	}
	return results
}

func testMap() {
	mapArray := []interface{}{
		map[string]interface{}{"key1": 1},
		map[string]interface{}{"key2": "2"},
	}

	fmt.Println(mapArray)

	fmt.Println(reflect.TypeOf(mapArray))

	var newMapArray []map[string]interface{}
	/*	newMapArray = mapArray*/
	/*	newMapArray = []map[string]interface{}(mapArray)*/
	for _, v := range mapArray {
		mapValue, ok := v.(map[string]interface{})
		if ok {
			newMapArray = append(newMapArray, mapValue)
		}
	}

	fmt.Println(newMapArray)
	fmt.Println(reflect.TypeOf(newMapArray))

}

// Flatten takes a map and returns a new one where nested maps are replaced
// by dot-delimited keys.
func Flatten(m map[string]interface{}) map[string]interface{} {
	o := make(map[string]interface{})
	for k, v := range m {
		switch child := v.(type) {
		case map[string]interface{}:
			nm := Flatten(child)
			for nk, nv := range nm {
				o[k+"_"+nk] = nv
			}
		default:
			o[k] = v
		}
	}
	return o
}

// LoadMXJ loads a dataset from a XML/JSON source.
// - MXJ allows to decode / encode XML or JSON to/from map[string]interface{};
//   extract values with dot-notation paths and wildcards.
// - Forked from `github.com/clbanning/mxj`
func LoadMXJ(jsonContent []byte) (*Dataset, error) {

	mxj.JsonUseNumber = true
	mv, err := mxj.NewMapJsonArray(jsonContent)
	if err != nil {
		return nil, err
	}

	// mxj.LeafUseDotNotation()
	// fmt.Println("TypeOf=", reflect.TypeOf(mv))

	var input []map[string]interface{}
	for _, v := range mv {
		// fmt.Println("k=", k, "TypeOf.v=", reflect.TypeOf(v))
		mapValue, ok := v.([]map[string]interface{})
		// fmt.Println("ok=", ok, "mapValue=", mapValue)
		if ok {
			for _, m := range mapValue {
				mx := Flatten(m)
				input = append(input, mx)
			}
		}
	}

	return internalLoadFromDict(input)
}

// func LoadMXJ(jsonContent []byte) (*Dataset, error)   { return nil, nil }
func LoadGJSON(jsonContent []byte) (*Dataset, error) { return nil, nil }

/*
// LoadGSON loads a dataset from a JSON source.
// - GJSON package allows to get values from a json document
//   with features such as one line retrieval, dot notation paths, iteration, and parsing json lines.
// - Forked from `github.com/tidwall/gjson`
func LoadGJSON(jsonContent []byte) (*Dataset, error) {

	// handle a goofy case ...
	//if jsonContent[0] == '[' {
	//	jsonContent = []byte(`{"array":` + string(jsonContent) + `}`)
	//}

	// if !gjson.Valid(string(jsonContent)) {
	// 	return nil, errors.New("invalid json")
	// }

	var input []map[string]interface{}
	if err := gjson.Unmarshal(jsonContent, &input); err != nil {
		return nil, err
	}


	//	input, ok := gjson.Parse(string(jsonContent)).Value().(map[string]interface{})
	//	if !ok {
	//		log.Fatalln("error, could not unmarshal to a map[string]interface{}")
	//		// not a map
	//	}
	//	pp.Println("map[string]interface{}=", input)

	// results := gjson.GetMany(json, "name.first", "name.last", "age")

	// var input []map[string]interface{}
	// results := gjson.GetMany(json, "name.first", "name.last", "age")
	// input = gjson.GetBytes(jsonContent, "").Value().([]map[string]interface{})

	result := gjson.ParseBytes(jsonContent) //.(map[string]interface{})
	pp.Println("Map=", result)

	// if !gjson.Valid(string(jsonContent)) {
	// 	return nil, errors.New("invalid json")
	// }

	// []map[string]interface{}
	// pp.Println("Array=", result.Array())
	// pp.Println("Map=", result.Map())

	// return nil, ErrUnmarshallingJsonWithGson
	return internalLoadFromDict(input)
}

*/

// LoadDatabookJSON loads a Databook from a JSON source.
func LoadDatabookJSON(jsonContent []byte) (*Databook, error) {
	var input []map[string]interface{}
	var internalInput []map[string]interface{}
	if err := json.Unmarshal(jsonContent, &input); err != nil {
		return nil, err
	}

	db := NewDatabook()
	for _, d := range input {
		b, err := json.MarshalIndent(d["data"], "", "\t")
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(b, &internalInput); err != nil {
			return nil, err
		}

		if ds, err := internalLoadFromDict(internalInput); err == nil {
			db.AddSheet(d["title"].(string), ds)
		} else {
			return nil, err
		}
	}

	return db, nil
}

// JSON returns a JSON representation of the Dataset as an Export.
func (d *Dataset) JSON() (*Export, error) {
	back := d.Dict()

	b, err := json.MarshalIndent(back, "", "\t")
	if err != nil {
		return nil, err
	}
	return newExportFromBytes(b), nil
}

// JSON returns a JSON representation of the Databook as an Export.
func (d *Databook) JSON() (*Export, error) {
	b := newBuffer()
	b.WriteString("[")
	for _, s := range d.sheets {
		b.WriteString("{\"title\": \"" + s.title + "\", \"data\": ")
		js, err := s.dataset.JSON()
		if err != nil {
			return nil, err
		}
		b.Write(js.Bytes())
		b.WriteString("},")
	}
	by := b.Bytes()
	by[len(by)-1] = ']'
	return newExportFromBytes(by), nil
}

func JsonMarshal(t interface{}) ([]byte, error) {
	buff := &bytes.Buffer{}
	enc := json.NewEncoder(buff)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	err := enc.Encode(t)
	return buff.Bytes(), err
}

func JsonMarshalIndent(t interface{}, prefix, indent string) ([]byte, error) {
	b, err := JsonMarshal(t)
	if err != nil {
		return b, err
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, prefix, indent)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
