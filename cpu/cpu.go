package cpu

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"os"
	"os/exec"
)

// CPUSpeed struct to hold the average and max clock speeds
type CPUSpeed struct {
	Avg int64
	Max int64
}

/// GetCPUClockSpeeds retrieves the average and maximum CPU clock speeds
func GetCPUClockSpeeds() (*CPUSpeed, error) {
    paths, err := filepath.Glob("/sys/devices/system/cpu/cpu*/cpufreq/scaling_cur_freq")
    if err != nil || len(paths) == 0 {
        return &CPUSpeed{Avg: -1, Max: -1}, nil // Return -1 if no data
    }

    var totalFreq, maxFreq int64
    coreCount := len(paths)

    for _, path := range paths {
        data, err := ioutil.ReadFile(path)
        if err != nil {
            return &CPUSpeed{Avg: -1, Max: -1}, nil // Return -1 on error
        }

        freqKHz, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
        if err != nil {
            return &CPUSpeed{Avg: -1, Max: -1}, nil // Return -1 on error
        }

        totalFreq += freqKHz
        if freqKHz > maxFreq {
            maxFreq = freqKHz
        }
    }

    avgFreq := totalFreq / int64(coreCount) / 1000 // in MHz
    maxFreqMHz := maxFreq / 1000                   // in MHz

    return &CPUSpeed{
        Avg: avgFreq,
        Max: maxFreqMHz,
    }, nil
}

// GetCPUTemperature retrieves the CPU temperature
func GetCPUTemperature() (string, error) {
	tempPath := "/sys/class/thermal/thermal_zone0/temp"
	if _, err := os.Stat(tempPath); err == nil {
		data, err := ioutil.ReadFile(tempPath)
		if err == nil {
			temp, err := strconv.Atoi(strings.TrimSpace(string(data)))
			if err == nil {
				return fmt.Sprintf("%.1f°C", float64(temp)/1000.0), nil
			}
		}
	}

	// Fallback to using lm-sensors if thermal_zone0 is not available
	cmd := exec.Command("sensors")
	out, err := cmd.Output()
	if err != nil {
		return "0°C", nil // Return "0°C" if sensors command fails
	}

	output := string(out)
	var temps []float64
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)

		// Check for AMD's Tccd lines or Intel's Core lines
		if len(fields) > 1 && (strings.HasPrefix(fields[0], "Tccd") || strings.HasPrefix(fields[0], "Core")) {
			// Extract the temperature value (second column, e.g., "+31.0°C")
			tempStr := strings.TrimSuffix(fields[1], "°C")
			tempStr = strings.TrimPrefix(tempStr, "+") // Remove '+' if present
			temp, err := strconv.ParseFloat(tempStr, 64)
			if err == nil {
				temps = append(temps, temp)
			}
		}
	}

	// Average temperature if multiple values found
	if len(temps) > 0 {
		totalTemp := 0.0
		for _, temp := range temps {
			totalTemp += temp
		}
		return fmt.Sprintf("%.1f°C", totalTemp/float64(len(temps))), nil
	}

	return "0°C", nil // Return "0°C" if no temperature is found
}
