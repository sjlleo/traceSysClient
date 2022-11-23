package main

import (
	"github.com/sjlleo/traceSysClient/config"
	"github.com/sjlleo/traceSysClient/traceService"
)

func main() {
	config.InitConfig()
	traceService.StartService()
}
