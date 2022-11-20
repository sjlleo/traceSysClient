package config

import (
	"log"
	"testing"

	"github.com/spf13/viper"
)

func TestConfig(t *testing.T) {
	InitConfig()
	log.Println(viper.Get("token"))
}
