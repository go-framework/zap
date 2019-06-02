package lumberjack_test

import (
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/go-framework/zap"
	"github.com/go-framework/zap/syncer"
	"github.com/go-framework/zap/syncer/lumberjack"
)

func Test_LumberjackLoggerSyncer(t *testing.T) {
	config := zap.GetDebugConfig()

	config.AddSyncerWrite(&syncer.Write{
		Name:   lumberjack.Name,
		Config: lumberjack.New("Test_LumberjackLoggerSyncer.log"),
	})

	t.Log("config",config)

	logger := config.NewZapLogger()

	logger.Info("Test_LumberjackLoggerSyncer")
}


func TestLumberjackLogger_UnmarshalYAML(t *testing.T) {
	filename := "config.yaml"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(filename, string(data))

	config := &zap.Config{}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("config", config)
}