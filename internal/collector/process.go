package collector

import (
    "github.com/shirou/gopsutil/v3/process"
    "time"
    "ddgo/internal/db"
)

func (sc *SystemCollector) collectProcessMetrics() {
    processes, err := process.Processes()
    if err != nil {
        return
    }

    for _, p := range processes {
        name, err := p.Name()
        if err != nil {
            continue
        }

        cpu, err := p.CPUPercent()
        if err != nil {
            continue
        }

        mem, err := p.MemoryPercent()
        if err != nil {
            continue
        }

        tags := make(map[string]string)
        for k, v := range sc.hostInfo {
            tags[k] = v
        }
        tags["process_name"] = name
        tags["pid"] = fmt.Sprintf("%d", p.Pid)

        metrics := []db.Metric{
            {
                Name:      "system.process.cpu",
                Value:     cpu,
                Tags:      tags,
                Timestamp: time.Now(),
            },
            {
                Name:      "system.process.memory",
                Value:     float64(mem),
                Tags:      tags,
                Timestamp: time.Now(),
            },
        }

        for _, metric := range metrics {
            if err := sc.db.SaveMetric(&metric); err != nil {
                log.Printf("Error saving process metric: %v", err)
            }
        }
    }
}