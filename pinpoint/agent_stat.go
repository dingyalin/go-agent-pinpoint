package pinpoint

import (
	"os"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"

	pinpoint "github.com/dingyalin/pinpoint-go-agent/thrift/dto/pinpoint"
)

func getTAgentStat(agentID string, startTime int64, collectInterval int64) *pinpoint.TAgentStat {
	timestamp := time.Now().UnixNano() / 1e6

	// process
	appCPULoad, appMEMUsage := getProcessStat()

	// cpu
	sysCPULoad := getSysCPULoad()
	cpuLoad := &pinpoint.TCpuLoad{
		JvmCpuLoad:    &appCPULoad,
		SystemCpuLoad: &sysCPULoad,
	}

	// mem
	memUsed, memTotal := getMemStat()
	gc := &pinpoint.TJvmGc{
		Type:                 0,
		JvmMemoryHeapUsed:    appMEMUsage,
		JvmMemoryHeapMax:     memTotal,
		JvmMemoryNonHeapUsed: memUsed,
		JvmMemoryNonHeapMax:  memTotal,
		JvmGcOldCount:        -1,
		JvmGcOldTime:         -1,
	}

	tagentStat := &pinpoint.TAgentStat{
		AgentId:         &agentID,
		StartTimestamp:  &startTime,
		Timestamp:       &timestamp,
		CollectInterval: &collectInterval,
		Gc:              gc,
		CpuLoad:         cpuLoad,
		Transaction:     nil, // nil
		ActiveTrace:     nil, // nil
		Metadata:        nil, // nil
	}

	return tagentStat
}

func getSysCPULoad() (sysCPULoad float64) {
	cpuPercents, err := cpu.Percent(0, false)
	if err != nil || len(cpuPercents) == 0 {
		return

	}

	sysCPULoad = cpuPercents[0] / float64(100)
	return
}

func getMemStat() (memUsed, memTotal int64) {
	mem, err := mem.VirtualMemory()
	if err != nil {
		return
	}

	memUsed = int64(mem.Used)
	memTotal = int64(mem.Total)
	return
}

func getProcessStat() (appCPULoad float64, appMEMUsage int64) {
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return
	}

	// cpu
	cpuPercent, err := proc.CPUPercent()
	if err != nil {
		return
	}
	appCPULoad = cpuPercent / float64(100)

	// mem
	mem, err := proc.MemoryInfo()
	appMEMUsage = int64(mem.RSS)

	return
}
