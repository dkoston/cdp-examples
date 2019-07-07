package logger

import (
	"github.com/op/go-logging"
	"os"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	// Ensure that if we set LOG_LEVEL=debug that our Logger instance reflects that
	err := os.Setenv("LOG_LEVEL", "debug")
	if err != nil {
		t.Error(err)
	}

	// Ensure we get *logging.Logger back
	log := Get()
	typeOfLog := reflect.TypeOf(log).String()
	expectedType := "*logging.Logger"
	if  typeOfLog !=  expectedType {
		t.Errorf("Failed to get expected type (%s). Got (%s)", expectedType, typeOfLog)
	}

	// Ensure our default Logger instance is using log level info
	if ! log.IsEnabledFor(logging.INFO) {
		t.Errorf("Expected log level INFO to be enabled but it IS NOT")
	}

	if ! log.IsEnabledFor(logging.DEBUG) {
		t.Errorf("Expected log level DEBUG to be enabled but it IS NOT")
	}
}
