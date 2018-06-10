package tablib

import (
	"bytes"
	"fmt"

	toml "github.com/BurntSushi/toml"
)

var ErrorOnUnmatchedKeys bool = false

// LoadTOML loads a dataset from a TOML source.
func LoadTOML(tomlContent []byte) (*Dataset, error) {
	var input []map[string]interface{}

	err := unmarshalToml(tomlContent, &input)
	if err != nil {
		return nil, err
	}

	return internalLoadFromDict(input)
}

// LoadDatabookTOML loads a Databook from a TOML source.
func LoadDatabookTOML(tomlContent []byte) (*Databook, error) {
	var input []map[string]interface{}
	var internalInput []map[string]interface{}

	err := unmarshalToml(tomlContent, &input)
	if err != nil {
		return nil, err
	}

	db := NewDatabook()
	for _, d := range input {
		/*
			 		b, err := yaml.Marshal(d["data"])
					if err != nil {
						return nil, err
					}
					if err := yaml.Unmarshal(b, &internalInput); err != nil {
						return nil, err
					}
		*/
		var b bytes.Buffer
		if dataBytes, err := getBytes(d["data"]); err == nil {
			if err := toml.NewEncoder(&b).Encode(dataBytes); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}

		if dataBytes, err := getBytes(d["data"]); err == nil {
			if err := unmarshalToml(dataBytes, &internalInput); err != nil {
				return nil, err
			}
		} else {
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

// TOML returns a TOML representation of the Dataset as an Export.
func (d *Dataset) TOML() (*Export, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	back := d.Dict()

	var b bytes.Buffer
	if err := toml.NewEncoder(&b).Encode(back); err != nil {
		return nil, err
	}

	return newExportFromBytes(b.Bytes()), nil
}

// TOML returns a TOML representation of the Databook as an Export.
func (d *Databook) TOML() (*Export, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	t := make([]map[string]interface{}, len(d.sheets))
	i := 0
	for _, s := range d.sheets {
		t[i] = make(map[string]interface{})
		t[i]["title"] = s.title
		t[i]["data"] = s.dataset.Dict()
		i++
	}

	var b bytes.Buffer
	if err := toml.NewEncoder(&b).Encode(t); err != nil {
		return nil, err
	}

	// return []byte(dataBytes.String()), nil

	return newExportFromBytes(b.Bytes()), nil
}

/*
func marshalToml(input []byte, output interface{}, errorOnUnmatchedKeys bool) error {
	var dataBytes bytes.Buffer
	if err := toml.NewEncoder(&dataBytes).Encode(t); err != nil {
		return nil, err
	}
}
*/

func unmarshalToml(input []byte, output interface{}) error {
	metadata, err := toml.Decode(string(input), output)
	if err == nil && len(metadata.Undecoded()) > 0 && ErrorOnUnmatchedKeys {
		return &UnmatchedTomlKeysError{Keys: metadata.Undecoded()}
	}
	return err
}

// UnmatchedTomlKeysError errors are returned by the Load function when
// ErrorOnUnmatchedKeys is set to true and there are unmatched keys in the input
// toml config file. The string returned by Error() contains the names of the
// missing keys.
type UnmatchedTomlKeysError struct {
	Keys []toml.Key
}

func (e *UnmatchedTomlKeysError) Error() string {
	return fmt.Sprintf("There are keys in the config file that do not match any field in the given struct: %v", e.Keys)
}
