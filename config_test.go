package zap

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestConfig_Parse(t *testing.T) {

	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))

	config := &Config{}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(config)
}

func TestConfig_MarshalJSON(t *testing.T) {

	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))

	config := &Config{}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(config)

	raw, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(raw))
}

func TestConfig_UnmarshalJSON(t *testing.T) {

	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))

	config := &Config{}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(config)

	raw, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	config2 := &Config{}

	err = json.Unmarshal(raw, config2)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(config2)
}
