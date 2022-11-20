package traceService

import (
	"testing"

	"github.com/sjlleo/traceSysClient/config"
)

func TestContab(t *testing.T) {
	config.InitConfig()
	StartService()
	// s.StartAsync()
}
