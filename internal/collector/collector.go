package collector 

import (
	"log"
	"sync"
	"time"
	"ddgo/internal/db"
	"github.com/shirou/gopsutil/v3/host"
)

type SystemCollector struct {
	db *db.dbinterval 
	interval time.Duration
	hostInfo map[string]string
	stopChan chan struct{}
	waitGroup sync.WaitGroup
}

func NewSystemCollector(database *db.DB, interval time.Duration) (*SystemCollector, error) {
	hostInfo, err := host.Info()

	if err != nill {
		return nil, err
	}

	return &SystemCollector{
		db: database,
		interval: interval,
		hostInfo: map[string]string {
			"hostname": hostInfo.Hostname,
			"os": hostInfo.OS,
			"platform": hostInfo.Platform,
		}
		stopChan: make(chan struct{}),
	}, nil
}

func (sc *SystemCollector) Start() {
	sc.WaitGroup.Add(1)
	go sc.collect()
}

func (sc *SystemCollector) collect() {
	defer sc.WaitGroup.Done()
	ticker := time.NewTicker(sc.interval)
	defer ticker.Stop()

	for {
		select {
		case <- ticker.C:
			sc.CollectMetrics()
		case <- sc.stopChan:
			return
		}
	}
}

func (sc *SystemCollector) collectMetrics() {
	sc.collectCPUMetrics()
	sc.collectMemoryMetrics()
	sc.collectDiskMetrics()
	sc.collectNetworkMetrics()
	sc.collectProcessMetrics()
}