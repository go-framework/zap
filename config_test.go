package zap

import (
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestConfig_UnmarshalYAML(t *testing.T) {
	filename := "config.yaml"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(filename, string(data))

	config := &Config{}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("config", config)
}
