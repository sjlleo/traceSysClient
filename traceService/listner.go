package traceService

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sjlleo/traceSysClient/fetchService"
	"github.com/sjlleo/traceSysClient/model"
	"github.com/sjlleo/traceSysClient/trace"
)

type TargetIP string

type Task struct {
	TaskIP          string
	NodeId          uint
	TaskId          uint
	IntervalSeconds int
	Scheduler       *gocron.Scheduler
	TraceConfig     *trace.Config
	ResultRWLock    *sync.RWMutex // 读写互斥锁
	TraceResult     []*trace.Result
}

const ReportCycle int = 20

var ActiveSchedulers map[TargetIP]Task

func init() {
	ActiveSchedulers = make(map[TargetIP]Task)
}

func ConfigInit() *trace.Config {
	defaultConfig := &trace.Config{
		BeginHop:         1,
		MaxHops:          30,
		NumMeasurements:  1,
		ParallelRequests: 18,
		Timeout:          1000 * time.Millisecond,
	}
	return defaultConfig
}

func (t *Task) UpdateScheduler(listptr *model.TraceList) {
	log.Println("traceService - UpdateScheduler")
	for _, item := range listptr.Task {
		if _, ok := ActiveSchedulers[TargetIP(item.IP)]; !ok {
			// 未开启的创建新的 Schedulers
			ActiveSchedulers[TargetIP(item.IP)] = Task{
				TaskId:          uint(item.TaskID),
				TaskIP:          item.IP,
				NodeId:          item.NodeID,
				IntervalSeconds: item.Interval,
				Scheduler:       gocron.NewScheduler(time.UTC),
				TraceConfig:     ConfigInit(),
				ResultRWLock:    &sync.RWMutex{},
			}
			// 获得 Map 对应 Key 的地址
			s := ActiveSchedulers[TargetIP(item.IP)]

			// 初始化配置
			dstIP := net.ParseIP(item.IP)
			s.TraceConfig.DestIP = dstIP
			s.TraceConfig.Method = trace.TraceMethod(item.Method)

			s.Scheduler.Every(item.Interval).Seconds().Do(s.GoTrace)
			// 每 6 个单位间隔汇报一次
			s.Scheduler.Every(item.Interval * ReportCycle).Seconds().Do(s.ResultReport)
			s.Scheduler.StartAsync()
		}
	}
}

func (t *Task) CleanUpScheduler(listptr *model.TraceList) {
	log.Println("traceService - CleanUpScheduler")

	for taskIP, activeTask := range ActiveSchedulers {
		var taskShouldDelete bool = true
		for _, pendingTask := range listptr.Task {
			if TargetIP(pendingTask.IP) == taskIP {
				taskShouldDelete = false
				break
			}
		}

		if taskShouldDelete {
			activeTask.Scheduler.Stop()
			delete(ActiveSchedulers, taskIP)
		}
	}
}

func (t *Task) CleanUpResult() {
	t.ResultRWLock.Lock()
	defer t.ResultRWLock.Unlock()

	t.TraceResult = t.TraceResult[len(t.TraceResult):]
}

func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

func (t *Task) ResultReport() {
	ReportMap := make(map[int]*model.HopReport)
	for ttl := 0; ttl < 30; ttl++ {
		var minLatency float64 = 9999
		var maxLatency float64
		var avgLatency float64
		var count int
		for _, j := range t.TraceResult {

			// 判断是否会溢出
			if ttl >= len(j.Hops) {
				break
			}

			// 判断Key是否初始化
			if ReportMap[ttl] == nil {
				ReportMap[ttl] = &model.HopReport{}
			}

			h := j.Hops[ttl][0]
			if h.Error == nil {
				ReportMap[ttl].IPList = append(ReportMap[ttl].IPList, h.Address.String())
				rtt := float64(h.RTT.Microseconds()) / 1000
				avgLatency += rtt
				if rtt > maxLatency {
					maxLatency = rtt
				}
				if rtt < minLatency {
					minLatency = rtt
				}
				count++
			}
		}
		if count != 0 {
			ReportMap[ttl].MaxLatency = maxLatency
			ReportMap[ttl].MinLatency = minLatency
			ReportMap[ttl].AvgLatency = avgLatency / float64(count)
			ReportMap[ttl].IPList = RemoveRepeatedElement(ReportMap[ttl].IPList)
			ReportMap[ttl].PacketLoss = (float64(len(t.TraceResult)) - float64(count)) / float64(len(t.TraceResult))
			// 防止偶尔因为路由跟踪线程延时导致的多记
			if ReportMap[ttl].PacketLoss < 0 {
				ReportMap[ttl].PacketLoss = 0
			}
		}
	}
	res := &model.Report{
		Data:     ReportMap,
		Interval: t.IntervalSeconds,
		NodeID:   t.NodeId,
		TaskID:   t.TaskId,
		Method:   int(t.TraceConfig.Method),
	}
	fetchService.PostResult(res)
	t.CleanUpResult()
}

func (t *Task) getConfig() {
	// 获取配置
	listptr, err := fetchService.FetchTraceList()
	if err != nil {
		return
	}
	go t.UpdateScheduler(listptr)
	go t.CleanUpScheduler(listptr)
}

func NewTask() *Task {
	return &Task{
		Scheduler: gocron.NewScheduler(time.UTC),
	}
}

func StartService() {
	t := NewTask()
	log.Println("traceService - StartService")
	t.Scheduler.Every("10s").Do(t.getConfig)
	t.Scheduler.StartBlocking()
}

func (t *Task) GoTrace() {
	var res *trace.Result
	switch t.TraceConfig.Method {
	case trace.ICMP:
		tracer := &trace.ICMPTracer{Config: *t.TraceConfig}
		res, _ = tracer.Execute()

	case trace.TCP:
		t.TraceConfig.DestPort = 443
		tracer := &trace.TCPTracer{Config: *t.TraceConfig}
		res, _ = tracer.Execute()

	case trace.UDP:
		t.TraceConfig.DestPort = 53
		tracer := &trace.UDPTracer{Config: *t.TraceConfig}
		res, _ = tracer.Execute()

	}
	t.ResultRWLock.Lock()
	defer t.ResultRWLock.Unlock()
	// log.Println(res)
	t.TraceResult = append(t.TraceResult, res)
}
