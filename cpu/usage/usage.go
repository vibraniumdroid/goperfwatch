package usage

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

// CPUUsage struct to hold current CPU usage percentage
type CPUUsage struct {
	Percent int
}

// readCPUTimes reads the CPU times from /proc/stat
func readCPUTimes() (idle, total int64, err error) {
	data, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read /proc/stat: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return 0, 0, fmt.Errorf("failed to parse /proc/stat data")
	}

	fields := strings.Fields(lines[0]) // First line should read "cpu ..."
	if len(fields) < 5 {
		return 0, 0, fmt.Errorf("unexpected format in /proc/stat")
	}

	// Parse CPU times from /proc/stat
	var cpuTimes []int64
	for _, field := range fields[1:] {
		value, err := strconv.ParseInt(field, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse CPU time: %v", err)
		}
		cpuTimes = append(cpuTimes, value)
	}

	// Calculate total and idle times
	idle = cpuTimes[3] // idle time is fourth field in /proc/stat
	for _, time := range cpuTimes {
		total += time
	}

	return idle, total, nil
}

// CalculateCPUUsage calculates the CPU usage over a short interval
func CalculateCPUUsage() (*CPUUsage, error) {
	idle1, total1, err := readCPUTimes()
	if err != nil {
		return &CPUUsage{Percent: -1}, nil // Return -1 if error
	}

	time.Sleep(250 * time.Millisecond) // wait before polling again 

	idle2, total2, err := readCPUTimes()
	if err != nil {
		return &CPUUsage{Percent: -1}, nil // Return -1 if error
	}

	idleDelta := idle2 - idle1
	totalDelta := total2 - total1

	usagePercent := 100 * (totalDelta - idleDelta) / totalDelta

	return &CPUUsage{Percent: int(usagePercent)}, nil
}
