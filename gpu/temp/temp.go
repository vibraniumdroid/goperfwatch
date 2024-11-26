package temp

import (
	"os/exec"
	"strings"
)

// GPUTemperature holds the GPU temperature reading
type GPUTemperature struct {
	EdgeTemp string
}

// GetGPUTemperature retrieves the GPU edge temperature from lm-sensors
func GetGPUTemperature() (*GPUTemperature, error) {
    cmd := exec.Command("sensors")
    output, err := cmd.Output()
    if err != nil {
        return &GPUTemperature{EdgeTemp: "-1°C"}, nil // Return -1°C if command fails
    }

    lines := strings.Split(string(output), "\n")
    for _, line := range lines {
        if strings.Contains(line, "edge:") {
            fields := strings.Fields(line)
            if len(fields) >= 2 {
                temp := strings.TrimPrefix(fields[1], "+")
                return &GPUTemperature{EdgeTemp: temp}, nil
            }
        }
    }

    return &GPUTemperature{EdgeTemp: "-1°C"}, nil // Return -1°C if no data found
}
