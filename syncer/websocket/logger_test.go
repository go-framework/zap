package websocket_test

import (
	"context"
	"io/ioutil"
	"log"
	"strings"
	"testing"
	"time"

	zap2 "go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/go-framework/zap"
	"github.com/go-framework/zap/syncer"
	"github.com/go-framework/zap/syncer/websocket"
)

func TestNew(t *testing.T) {
	config := zap.GetDebugConfig()

	u := "ws://127.0.0.1:9008/realtime/logger"

	write := websocket.New(u)

	defer write.Close()

	config.AddSyncerWrite(&syncer.Write{
		Name:   websocket.Name,
		Config: write,
	})

	// config.Console = false

	t.Log("config", config)

	WebSocketLoggerSyncer(t, config)
}

func WebSocketLoggerSyncer(t *testing.T, config *zap.Config) {

	logger := config.NewZapLogger()
	defer logger.Sync()

	logger.Info("Test_WebSocketLoggerSyncer")

	ctx, _ := context.WithTimeout(context.Background(), time.Second*600)

	duration := time.Second
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	count := 0

	const times = 1000
	logger.WithOptions()
	logger.Info("envelope", zap2.String("long", strings.Repeat("0123456789", times)))

	for {
		select {
		case <-ticker.C:
			count++
			// if count > 5 {
			// 	ticker.Stop()
			// 	break
			// }
			// duration *= 2
			// if duration > 10*time.Second {
			// 	duration = time.Second
			// }
			// ticker = time.NewTicker(duration)
			logger.Info("loop", zap2.Duration("duration", duration), zap2.Int("count", count), zap2.String("match", strings.Repeat(string('a'+rune(count%26)),5)))

		case <-ctx.Done():
			log.Println("exit")
			logger.Info("exit", zap2.Error(ctx.Err()))
			return
		}
	}
}

func TestWebSocketLogger_UnmarshalYAML(t *testing.T) {
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

func TestWebSocketLogger_ConfigFile(t *testing.T) {
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

	WebSocketLoggerSyncer(t, config)
}
