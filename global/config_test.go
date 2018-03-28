package global

import (
	"testing"
	"github.com/spf13/viper"
)

func TestConfigLoaded(t *testing.T) {
	err := LoadConfig()
	if err != nil {
		t.Fatalf("Fatal error config file: %s \n", err)
	}

	if !viper.IsSet("server") {
		t.Error("config file expected key server but it isn't exists")
	}

	if !viper.IsSet("database") {
		t.Error("config file expected key database but it isn't exists")
	}
}
