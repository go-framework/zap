package syncer

import (
	"errors"
	"fmt"
	"io"

	"github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
)

// Global enabled Writers.
var gWriters = make(map[string]io.Writer)

// Register Writer.
// Writer should tag as `json:",inline" yaml:",inline" mapstructure:",squash"` format.
func RegisterWriter(name string, writer io.Writer) {
	gWriters[name] = writer
}

// Defined writers for Get Writer from config.
type Write struct {
	// Write name.
	Name string `json:"name" yaml:"name"`
	// Write config which implement Writer.
	Config io.Writer `json:"config" yaml:"config"`
}

// Get Writer.
func (this *Write) GetWriter() io.Writer {
	return this.Config
}

// unmarshal map[string]interface{}) data, get the name and config is exist,
// the Writer which implement structure should be tag as `json:",inline" yaml:",inline" mapstructure:",squash"` format.
func (this *Write) unmarshal(data map[string]interface{}) error {
	// get write name.
	if name, ok := data["name"]; ok {
		switch v := name.(type) {
		case string:
			this.Name = v
		default:
			this.Name = fmt.Sprintf("%v", name)
		}
	} else {
		return errors.New("write should be have name filed")
	}

	// load Writer.
	if writer, ok := gWriters[this.Name]; ok {
		// if implement Cloner then new one.
		if face, ok := writer.(Cloner); ok {
			this.Config = face.Clone()
		} else {
			this.Config = writer
		}
	} else {
		return errors.New("not support write: " + this.Name)
	}

	// if have config filed then parse it.
	if config, ok := data["config"]; ok {
		err := mapstructure.Decode(config, this.Config)
		if err != nil {
			return err
		}
	}

	return nil
}

// Implement YAML Unmarshaler interface.
func (this *Write) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// new temp as map[string]interface{}.
	temp := make(map[string]interface{})

	// unmarshal temp as yaml.
	if err := unmarshal(temp); err != nil {
		return err
	}

	return this.unmarshal(temp)
}

// Implement JSON Unmarshaler interface.
func (this *Write) UnmarshalJSON(data []byte) error {

	// new temp as map[string]interface{}.
	temp := make(map[string]interface{})

	// unmarshal temp as json.
	err := jsoniter.Unmarshal(data, &temp)
	if err != nil {
		return err
	}

	return this.unmarshal(temp)
}
