package memory

import (
	"io/ioutil"
	"strconv"
	"strings"
)

// MemoryStats holds memory usage information
type MemoryStats struct {
	UsedGB  float64
	TotalGB float64
	Percent int
}

// GetMemoryUsage calculates memory usage in GB and percentage
func GetMemoryUsage() (*MemoryStats, error) {
	data, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return &MemoryStats{UsedGB: -1, TotalGB: -1, Percent: -1}, nil // Return -1 if error
	}

	memInfo := map[string]int64{}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSuffix(fields[0], ":")
		value, err := strconv.ParseInt(fields[1], 10, 64)
		if err == nil {
			memInfo[key] = value
		}
	}

	total := memInfo["MemTotal"]
	free := memInfo["MemFree"]
	buffers := memInfo["Buffers"]
	cached := memInfo["Cached"]

	used := total - free - buffers - cached
	usedGB := float64(used) / (1024 * 1024) // GB
	totalGB := float64(total) / (1024 * 1024) // GB

	percent := 0
	if total > 0 {
		percent = int((float64(used) / float64(total)) * 100)
	}

	return &MemoryStats{
		UsedGB:  usedGB,
		TotalGB: totalGB,
		Percent: percent,
	}, nil
}
