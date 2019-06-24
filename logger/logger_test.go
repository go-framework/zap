package logger

import (
	"testing"
)

func TestInfo(t *testing.T) {
	Info("info")
}

func TestFatal(t *testing.T) {
	Fatal("fatal")
}
