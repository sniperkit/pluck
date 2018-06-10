package config

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestNewParser_ok(t *testing.T) {
	configPath := "/tmp/ok.json"
	configContent := []byte(`{
    "version": 2,
    "name": "My lovely gateway",
    "port": 8080,
    "cache_ttl": "3600s",
    "timeout": "3s",
    "endpoints": [
        {
            "endpoint": "/github",
            "method": "GET",
            "extra_config" : {"user":"test","hits":6,"parents":["gomez","morticia"]},
            "backend": [
                {
                    "host": [
                        "https://api.github.com"
                    ],
                    "url_pattern": "/",
                    "whitelist": [
                        "authorizations_url",
                        "code_search_url"
                    ],
                    "extra_config" : {"user":"test","hits":6,"parents":["gomez","morticia"]}
                }
            ]
        },
        {
            "endpoint": "/supu",
            "method": "GET",
            "concurrent_calls": 3,
            "backend": [
                {
                    "host": [
                        "http://127.0.0.1:8080"
                    ],
                    "url_pattern": "/__debug/supu"
                }
            ]
        },
        {
            "endpoint": "/combination/{id}",
            "method": "GET",
            "backend": [
                {
                    "group": "first_post",
                    "host": [
                        "https://jsonplaceholder.typicode.com"
                    ],
                    "url_pattern": "/posts/{id}",
                    "blacklist": [
                        "userId"
                    ]
                },
                {
                    "host": [
                        "https://jsonplaceholder.typicode.com"
                    ],
                    "url_pattern": "/users/{id}",
                    "mapping": {
                        "email": "personal_email"
                    }
                }
            ]
        }
    ],
    "extra_config" : {"user":"test","hits":6,"parents":["gomez","morticia"]}
}`)
	if err := ioutil.WriteFile(configPath, configContent, 0644); err != nil {
		t.FailNow()
	}

	serviceConfig, err := NewParser().Parse(configPath)

	if err != nil {
		t.Error("Unexpected error. Got", err.Error())
	}
	testExtraConfig(serviceConfig.ExtraConfig, t)

	endpoint := serviceConfig.Endpoints[0]
	endpointExtraConfiguration := endpoint.ExtraConfig

	if endpointExtraConfiguration != nil {
		testExtraConfig(endpointExtraConfiguration, t)
	} else {
		t.Error("Extra config is not present in EndpointConfig")
	}

	backend := endpoint.Backend[0]
	backendExtraConfiguration := backend.ExtraConfig
	if backendExtraConfiguration != nil {
		testExtraConfig(backendExtraConfiguration, t)
	} else {
		t.Error("Extra config is not present in BackendConfig")
	}

	if err := os.Remove(configPath); err != nil {
		t.FailNow()
	}
}

func testExtraConfig(extraConfig map[string]interface{}, t *testing.T) {
	userVar := extraConfig["user"]
	if userVar != "test" {
		t.Error("User in extra config is not test")
	}
	parents := extraConfig["parents"].([]interface{})
	if parents[0] != "gomez" {
		t.Error("Parent 0 of user us not gomez")
	}
	if parents[1] != "morticia" {
		t.Error("Parent 1 of user us not morticia")
	}
}

func TestNewParser_unknownFile(t *testing.T) {
	_, err := NewParser().Parse("/nowhere/in/the/fs.json")
	if err == nil || strings.Index(err.Error(), "Fatal error config file:") != 0 {
		t.Error("Error expected. Got", err)
	}
}

func TestNewParser_readingError(t *testing.T) {
	wrongConfigPath := "/tmp/reading.json"
	wrongConfigContent := []byte("{hello\ngo\n")
	if err := ioutil.WriteFile(wrongConfigPath, wrongConfigContent, 0644); err != nil {
		t.FailNow()
	}

	expected := "Fatal error config file: While parsing config: invalid character 'h' looking for beginning of object key string"
	_, err := NewParser().Parse(wrongConfigPath)
	if err == nil || strings.Index(err.Error(), expected) != 0 {
		t.Error("Error expected. Got", err)
	}
	if err = os.Remove(wrongConfigPath); err != nil {
		t.FailNow()
	}
}

func TestNewParser_initError(t *testing.T) {
	wrongConfigPath := "/tmp/unmarshall.json"
	wrongConfigContent := []byte("{\"a\":42}")
	if err := ioutil.WriteFile(wrongConfigPath, wrongConfigContent, 0644); err != nil {
		t.FailNow()
	}

	_, err := NewParser().Parse(wrongConfigPath)
	if err == nil || strings.Index(err.Error(), "Unsupported version: 0") != 0 {
		t.Error("Error expected. Got", err)
	}
	if err = os.Remove(wrongConfigPath); err != nil {
		t.FailNow()
	}
}

func TestParserFunc(t *testing.T) {
	expected := ServiceConfig{Version: 42}
	result, err := ParserFunc(func(_ string) (ServiceConfig, error) { return expected, nil })("path/to/the/config/file")
	if err != nil {
		t.Error(err.Error())
	}
	if result.Version != expected.Version {
		t.Error("unexpected parsed config:", result)
	}
}
