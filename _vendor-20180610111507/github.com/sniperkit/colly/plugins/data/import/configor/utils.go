package configor

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"

	toml "github.com/BurntSushi/toml"
	ini "github.com/go-ini/ini"
	// yaml "gopkg.in/yaml.v2"
	yaml "github.com/go-yaml/yaml"

	pp "github.com/sniperkit/xutil/plugin/debug/pp"
)

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

func (configor *Configor) getENVPrefix(config interface{}) string {
	if configor.Config.ENVPrefix == "" {
		if prefix := os.Getenv("CONFIGOR_ENV_PREFIX"); prefix != "" {
			return prefix
		}
		return "Configor"
	}
	return configor.Config.ENVPrefix
}

func getConfigurationFileWithENVPrefix(file, env string) (string, error) {
	var (
		envFile string
		extname = path.Ext(file)
	)

	if extname == "" {
		envFile = fmt.Sprintf("%v.%v", file, env)
	} else {
		envFile = fmt.Sprintf("%v.%v%v", strings.TrimSuffix(file, extname), env, extname)
	}

	if fileInfo, err := os.Stat(envFile); err == nil && fileInfo.Mode().IsRegular() {
		return envFile, nil
	}
	return "", fmt.Errorf("failed to find file %v", file)
}

func (configor *Configor) getConfigurationFiles(files ...string) []string {
	var results []string

	if configor.Config.Debug || configor.Config.Verbose {
		fmt.Printf("Current environment: '%v'\n", configor.GetEnvironment())
	}

	for i := len(files) - 1; i >= 0; i-- {
		foundFile := false
		file := files[i]

		// check configuration
		if fileInfo, err := os.Stat(file); err == nil && fileInfo.Mode().IsRegular() {
			foundFile = true
			results = append(results, file)
		}

		// check configuration with env
		if file, err := getConfigurationFileWithENVPrefix(file, configor.GetEnvironment()); err == nil {
			foundFile = true
			results = append(results, file)
		}

		// check example configuration
		if !foundFile {
			if example, err := getConfigurationFileWithENVPrefix(file, "example"); err == nil {
				fmt.Printf("Failed to find configuration %v, using example file %v\n", file, example)
				results = append(results, example)
			} else {
				fmt.Printf("Failed to find configuration %v\n", file)
			}
		}
	}
	return results
}

func processFile(config interface{}, file string, errorOnUnmatchedKeys bool) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	switch {
	case strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml"):
		if errorOnUnmatchedKeys {
			fmt.Println("switch.yaml.UnmarshalStrict().errorOnUnmatchedKeys=", errorOnUnmatchedKeys)
			return yaml.UnmarshalStrict(data, config)
		}
		fmt.Println("switch.yaml.Unmarshal().errorOnUnmatchedKeys=", errorOnUnmatchedKeys)
		return yaml.Unmarshal(data, config)

	case strings.HasSuffix(file, ".toml"):
		return unmarshalToml(data, config, errorOnUnmatchedKeys)

	case strings.HasSuffix(file, ".json"):
		return json.Unmarshal(data, config)

	case strings.HasSuffix(file, ".xml"):
		return xml.Unmarshal(data, config)

	default:

		if err := unmarshalToml(data, config, errorOnUnmatchedKeys); err == nil {
			return nil
		} else if errUnmatchedKeys, ok := err.(*UnmatchedTomlKeysError); ok {
			return errUnmatchedKeys
		}

		if json.Unmarshal(data, config) == nil {
			return nil
		}

		var yamlError error
		if errorOnUnmatchedKeys {
			yamlError = yaml.UnmarshalStrict(data, config)
		} else {
			fmt.Println("default.yaml.UnmarshalStrict().errorOnUnmatchedKeys=", errorOnUnmatchedKeys)
			fmt.Println("default.yaml.Unmarshal().errorOnUnmatchedKeys=", errorOnUnmatchedKeys)
			yamlError = yaml.Unmarshal(data, config)
		}

		if yamlError == nil {
			return nil
		} else if yErr, ok := yamlError.(*yaml.TypeError); ok {
			return yErr
		}

		return errors.New("failed to decode config")
	}
}

// GetStringTomlKeys returns a string array of the names of the keys that are passed in as args
func GetStringTomlKeys(list []toml.Key) []string {
	arr := make([]string, len(list))

	for index, key := range list {
		arr[index] = key.String()
	}
	return arr
}

func unmarshalToml(data []byte, config interface{}, errorOnUnmatchedKeys bool) error {
	metadata, err := toml.Decode(string(data), config)
	if err == nil && len(metadata.Undecoded()) > 0 && errorOnUnmatchedKeys {
		return &UnmatchedTomlKeysError{Keys: metadata.Undecoded()}
	}
	return err
}

func getPrefixForStruct(prefixes []string, fieldStruct *reflect.StructField) []string {
	if fieldStruct.Anonymous && fieldStruct.Tag.Get("anonymous") == "true" {
		return prefixes
	}
	return append(prefixes, fieldStruct.Name)
}

func (configor *Configor) processTags(config interface{}, prefixes ...string) error {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	if configValue.Kind() != reflect.Struct {
		return errors.New("invalid config, should be struct")
	}

	configType := configValue.Type()
	for i := 0; i < configType.NumField(); i++ {
		var (
			envNames    []string
			fieldStruct = configType.Field(i)
			field       = configValue.Field(i)
			envName     = fieldStruct.Tag.Get("env") // read configuration from shell env
		)

		if !field.CanAddr() || !field.CanInterface() {
			continue
		}

		if envName == "" {
			envNames = append(envNames, strings.Join(append(prefixes, fieldStruct.Name), "_"))                  // Configor_DB_Name
			envNames = append(envNames, strings.ToUpper(strings.Join(append(prefixes, fieldStruct.Name), "_"))) // CONFIGOR_DB_NAME
		} else {
			envNames = []string{envName}
		}

		if configor.Config.Verbose {
			fmt.Printf("Trying to load struct `%v`'s field `%v` from env %v\n", configType.Name(), fieldStruct.Name, strings.Join(envNames, ", "))
		}

		// Load From Shell ENV
		for _, env := range envNames {
			if value := os.Getenv(env); value != "" {
				if configor.Config.Debug || configor.Config.Verbose {
					fmt.Printf("Loading configuration for struct `%v`'s field `%v` from env %v...\n", configType.Name(), fieldStruct.Name, env)
				}
				if err := yaml.Unmarshal([]byte(value), field.Addr().Interface()); err != nil {
					return err
				}
				break
			}
		}

		if isBlank := reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()); isBlank {
			// Set default configuration if blank
			if value := fieldStruct.Tag.Get("default"); value != "" {
				if err := yaml.Unmarshal([]byte(value), field.Addr().Interface()); err != nil {
					return err
				}
			} else if fieldStruct.Tag.Get("required") == "true" {
				// return error if it is required but blank
				return errors.New(fieldStruct.Name + " is required, but blank")
			}
		}

		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		if field.Kind() == reflect.Struct {
			if err := configor.processTags(field.Addr().Interface(), getPrefixForStruct(prefixes, &fieldStruct)...); err != nil {
				return err
			}
		}

		if field.Kind() == reflect.Slice {
			for i := 0; i < field.Len(); i++ {
				if reflect.Indirect(field.Index(i)).Kind() == reflect.Struct {
					if err := configor.processTags(field.Index(i).Addr().Interface(), append(getPrefixForStruct(prefixes, &fieldStruct), fmt.Sprint(i))...); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func encodeFile(config interface{}, node string, format string) ([]byte, error) {
	switch format {
	case "ini":

		pp.Println("encodeInterfaceToINI=", config)
		err := ini.MapTo(config, "./colly.ini")
		if err != nil {
			fmt.Println("error: ", err)
			os.Exit(1)
			return nil, err
		}

		// cfg := ini.Empty()

		// cfg.WriteToIndent(writer, "\t")
		// err = cfg.SaveToIndent("my.ini", "\t")
		/*
			var dataBytes bytes.Buffer
			// dataBytes := bytes.NewBuffer(nil)
			// configValue := reflect.ValueOf(config)
			outFile := ini.Empty()
			configValue := reflect.Indirect(reflect.ValueOf(*config))
			configValue2 := reflect.ValueOf(config)
			pp.Println(config)
			pp.Println(configValue)
			pp.Println(configValue2)
			err := ini.ReflectFrom(outFile, config)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			output, err := outFile.WriteTo(&dataBytes)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			fmt.Println(output)
			return []byte(dataBytes.String()), nil
		*/
	case "json":
		data, err := json.MarshalIndent(config, "", "\t")
		if err != nil {
			return nil, err
		}
		return data, nil

	case "toml":
		var dataBytes bytes.Buffer
		if err := toml.NewEncoder(&dataBytes).Encode(config); err != nil {
			return nil, err
		}
		// fmt.Println(dataBytes.String())
		return []byte(dataBytes.String()), nil

	case "xml":
		data, err := xml.MarshalIndent(config, "", "\t")
		if err != nil {
			return nil, err
		}
		return data, nil

	case "yaml":
		data, err := yaml.Marshal(config)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, errors.New("Unkown format to export")

}
