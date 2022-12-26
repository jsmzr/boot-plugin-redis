package redis

import (
	"testing"

	"github.com/spf13/viper"
)

func TestBase(t *testing.T) {
	p := RedisPlugin{}
	if p.Enabled() != defaultConfig["enabled"] {
		t.Fatalf("enabled should be %v", defaultConfig["enabled"])
	}
	if p.Order() != defaultConfig["order"] {
		t.Fatalf("order should be %v", defaultConfig["order"])
	}

}

func TestLoad(t *testing.T) {
	p := RedisPlugin{}
	viper.Set(configPrefix+"type", "error")
	if p.Load() == nil {
		t.Fatal("load error type")
	}

}
