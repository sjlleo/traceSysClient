package model

type TraceList struct {
	Task []TraceTask `json:"task"`
}

type TraceTask struct {
	TaskID     int    `json:"id"`
	Interval   int    `json:"interval"`
	Method     int    `json:"method"`
	NodeID     uint   `json:"nodeId"`
	TargetPort int    `json:"targetPort"`
	IP         string `json:"ip"`
}

type HopReport struct {
	IPList     []string `json:"ip_list"`
	PacketLoss float64  `json:"packetLoss"`
	MinLatency float64  `json:"min_latency"`
	MaxLatency float64  `json:"max_latency"`
	AvgLatency float64  `json:"avg_latency"`
}

type Report struct {
	Data     map[int]*HopReport `json:"data"`
	NodeID   uint               `json:"nodeId"`
	TaskID   uint               `json:"taskId"`
	Interval int                `json:"interval"`
	Token    string             `json:"token"`
	Method   int                `json:"method"`
}
