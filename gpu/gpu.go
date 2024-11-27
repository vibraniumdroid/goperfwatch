package gpu

import (
	"fmt"
	"os/exec"
	"strings"
)

// VRAMUsage holds the VRAM usage details
type VRAMUsage struct {
	UsedMB  int64
	TotalMB int64
}

// GetVRAMUsage retrieves the VRAM usage in MB
func GetVRAMUsage() (*VRAMUsage, error) {
	cmd := exec.Command("glxinfo")
	output, err := cmd.Output()
	if err != nil {
		return &VRAMUsage{UsedMB: -1, TotalMB: -1}, nil // Return -1 if command fails
	}

	var totalVRAM, availableVRAM int64
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "Dedicated video memory") {
			fmt.Sscanf(line, "    Dedicated video memory: %d MB", &totalVRAM)
		}
		if strings.Contains(line, "Currently available dedicated video memory") {
			fmt.Sscanf(line, "    Currently available dedicated video memory: %d MB", &availableVRAM)
		}
	}

	if totalVRAM == 0 || availableVRAM == 0 {
		return &VRAMUsage{UsedMB: -1, TotalMB: -1}, nil // Return -1 if data not found
	}

	usedVRAM := totalVRAM - availableVRAM

	return &VRAMUsage{
		UsedMB:  usedVRAM,
		TotalMB: totalVRAM,
	}, nil
}
