package collector

import (
    "github.com/shirou/gopsutil/v3/net"
    "time"
    "ddgo/internal/db"
)

func (sc *SystemCollector) collectNetworkMetrics() {
    interfaces, err := net.Interfaces()
    if err != nil {
        return
    }

    for _, iface := range interfaces {
        if len(iface.Addrs) == 0 {
            continue
        }

        ioCounters, err := net.IOCounters(true)
        if err != nil {
            continue
        }

        for _, io := range ioCounters {
            if io.Name == iface.Name {
                tags := make(map[string]string)
                for k, v := range sc.hostInfo {
                    tags[k] = v
                }
                tags["interface"] = iface.Name

                metrics := []db.Metric{
                    {
                        Name:      "system.net.bytes_sent",
                        Value:     float64(io.BytesSent),
                        Tags:      tags,
                        Timestamp: time.Now(),
                    },
                    {
                        Name:      "system.net.bytes_recv",
                        Value:     float64(io.BytesRecv),
                        Tags:      tags,
                        Timestamp: time.Now(),
                    },
                    {
                        Name:      "system.net.packets_sent",
                        Value:     float64(io.PacketsSent),
                        Tags:      tags,
                        Timestamp: time.Now(),
                    },
                    {
                        Name:      "system.net.packets_recv",
                        Value:     float64(io.PacketsRecv),
                        Tags:      tags,
                        Timestamp: time.Now(),
                    },
                }

                for _, metric := range metrics {
                    if err := sc.db.SaveMetric(&metric); err != nil {
                        log.Printf("Error saving network metric: %v", err)
                    }
                }
            }
        }
    }
}