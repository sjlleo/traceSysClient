package trace

import (
	"errors"
	"net"
	"sync"
	"time"
)

type TraceMethod int

var (
	ICMP                  TraceMethod = 0
	TCP                   TraceMethod = 1
	UDP                   TraceMethod = 2
	ErrInvalidMethod                  = errors.New("invalid method")
	ErrTracerouteExecuted             = errors.New("traceroute already executed")
	ErrHopLimitTimeout                = errors.New("hop timeout")
)

type Result struct {
	Hops [][]Hop
	lock sync.Mutex
}

type Config struct {
	Method            TraceMethod
	SrcAddr           string
	BeginHop          int
	MaxHops           int
	NumMeasurements   int
	ParallelRequests  int
	Timeout           time.Duration
	DestIP            net.IP
	DestPort          int
	Quic              bool
	RDns              bool
	CallBackInterface func(res *Result, ttl int)
}

type Hop struct {
	Success  bool
	Address  net.Addr
	Hostname string
	TTL      int
	RTT      time.Duration
	Error    error
}

func (s *Result) reduce(final int) {
	if final > 0 && final < len(s.Hops) {
		s.Hops = s.Hops[:final]
	}
}

func (s *Result) add(hop Hop) {
	s.lock.Lock()
	defer s.lock.Unlock()
	k := hop.TTL - 1
	for len(s.Hops) < hop.TTL {
		s.Hops = append(s.Hops, make([]Hop, 0))
	}
	s.Hops[k] = append(s.Hops[k], hop)

}
