package plucker

import (
	jsoniter "github.com/json-iterator/go"
	// xml "encoding/xml"
	// toml "github.com/BurntSushi/toml"
	// yaml "github.com/go-yaml/yaml"
	// ini "github.com/go-ini/ini"
	// cfg "github.com/sniperkit/colly/plugins/data/extract/text-plucker/config"
	// configor "github.com/sniperkit/colly/plugins/data/import/configor"
)

// faster than the default "encoding/json" package
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// ResultJSON returns the result, formatted as JSON.
// If their are no results, it returns an empty string.
func (p *Plucker) ResultJSON(indent ...bool) (string, error) {
	totalResults := 0
	for key := range p.result {
		b, _ := json.Marshal(p.result[key])
		totalResults += len(b)
	}
	if totalResults == len(p.result)*2 { // results == 2 because its just []
		// log.Error(errors.Wrap(err, errEmptyResults)
		return "", errEmptyResults
	}
	var err error
	var resultJSON []byte
	if len(indent) > 0 && indent[0] {
		resultJSON, err = json.MarshalIndent(p.result, "", "    ")
	} else {
		resultJSON, err = json.Marshal(p.result)
	}
	if err != nil {
		return "", err
		// log.Error(errors.Wrap(err, errMarshallingResults))
	}
	return string(resultJSON), nil
}

// Result returns the raw result
func (p *Plucker) Result() map[string]interface{} {
	return p.result
}

// LoadFromString will load a JSON/YAML/TOML/INI configuration file of units
// to pluck with specified parameters
func (p *Plucker) LoadFromString(tomlString string) (err error) {

	// ?! isYaml, isTOML, isINI, isXML
	// var conf configs
	// conf := &configs{}

	/*
		_, err = toml.Decode(tomlString, &conf)
		// log.Debugf("Loaded toml: %+v", conf)
		for i := range conf.Pluck {
			var c config.Config
			c.Activators = conf.Pluck[i].Activators
			c.Deactivator = conf.Pluck[i].Deactivator
			c.Finisher = conf.Pluck[i].Finisher
			c.Limit = conf.Pluck[i].Limit
			c.Name = conf.Pluck[i].Name
			c.Permanent = conf.Pluck[i].Permanent
			c.Sanitize = conf.Pluck[i].Sanitize
			c.Maximum = conf.Pluck[i].Maximum
			p.Add(c)
		}
	*/

	return
}

/*
func configParse(conf Configs) []cfg.Config {
	var pluckers Configs

	// pluckers := make()

	for i := range conf.Pluck {
		var c cfg.Config
		c.Activators = conf.Pluck[i].Activators
		c.Deactivator = conf.Pluck[i].Deactivator
		c.Finisher = conf.Pluck[i].Finisher
		c.Limit = conf.Pluck[i].Limit
		c.Name = conf.Pluck[i].Name
		c.Permanent = conf.Pluck[i].Permanent
		c.Sanitize = conf.Pluck[i].Sanitize
		c.Maximum = conf.Pluck[i].Maximum

		pluckers = append(pluckers, c)
	}
	return pluckers
}

func LoadJSON(content []byte) (configs []cfg.Config, err error) {
	var conf Configs
	err = json.Unmarshal(content, &conf)
	if err != nil {
		return nil, err
	}
	configs = configParse(conf)
	return
}

// loadToml will load a TOML configuration file of units to pluck with their specified parameters
func LoadTOML(content []byte) (configs []cfg.Config, err error) {
	var conf Configs
	_, err = toml.Decode(content, &conf)
	if err != nil {
		return nil, err
	}
	configs = configParse(conf)
	return
}

func LoadYAML(content []byte) (configs []cfg.Config, err error) {
	var conf Configs
	err = yaml.Unmarshal(content, &conf)
	if err != nil {
		return nil, err
	}
	configs = configParse(conf)
	return
}

func LoadXML(content []byte) (configs []cfg.Config, err error) {
	var conf Configs
	_, err = xml.Unmarshal(content, &conf)
	if err != nil {
		return nil, err
	}
	configs = configParse(conf)
	return
}
*/
