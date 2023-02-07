package trace

import (
	"log"
	"net"
	"testing"
	"time"
)

func ConfigInit() *Config {
	defaultConfig := &Config{
		BeginHop:         1,
		MaxHops:          30,
		NumMeasurements:  1,
		ParallelRequests: 18,
		Timeout:          1000 * time.Millisecond,
	}
	return defaultConfig
}

func TestICMPv4(t *testing.T) {
	c := ConfigInit()
	TargetIPStr := "101.95.52.34"
	c.DestIP = net.ParseIP(TargetIPStr)
	tracer := &ICMPTracer{Config: *c}
	res, _ := tracer.Execute()
	log.Println(res)
}

func TestTCPv4(t *testing.T) {
	c := ConfigInit()
	TargetIPStr := "8.8.8.8"
	c.DestIP = net.ParseIP(TargetIPStr)
	c.DestPort = 443
	tracer := &TCPTracer{Config: *c}
	res, _ := tracer.Execute()
	log.Println(res)
}

func TestUDPv4(t *testing.T) {
	c := ConfigInit()
	TargetIPStr := "8.8.8.8"
	c.DestIP = net.ParseIP(TargetIPStr)
	c.DestPort = 53
	tracer := &UDPTracer{Config: *c}
	res, _ := tracer.Execute()
	log.Println(res)
}
