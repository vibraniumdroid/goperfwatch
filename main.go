package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"goperfwatch/cpu"
	"goperfwatch/cpu/usage"
	"goperfwatch/memory"
	"goperfwatch/gpu"
	"goperfwatch/gpu/temp"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"os/exec"
)

func enableRawMode() (*exec.Cmd, *os.File) {
	cmd := exec.Command("stty", "-icanon", "-echo") // Turn off canonical mode and echo
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
	return cmd, os.Stdin
}

func disableRawMode(cmd *exec.Cmd) {
	cmd = exec.Command("stty", "icanon", "echo") // Restore terminal settings
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func main() {
	// Enable raw mode for keyboard input
	cmd, _ := enableRawMode()
	defer disableRawMode(cmd)

	// Initialize the termui library
	if err := termui.Init(); err != nil {
		fmt.Println("failed to initialize termui:", err)
		return
	}
	defer termui.Close()

	// Create widgets for CPU, memory, GPU VRAM, and GPU temperature
	clockSpeedWidget := widgets.NewParagraph()
	clockSpeedWidget.SetRect(0, 0, 50, 5)
	clockSpeedWidget.Title = "CPU Clock Speed and Temperature"

	usageGauge := widgets.NewGauge()
	usageGauge.SetRect(0, 5, 50, 8)
	usageGauge.Title = "CPU Usage"

	memoryGauge := widgets.NewGauge()
	memoryGauge.SetRect(0, 8, 50, 11)
	memoryGauge.Title = "Memory Usage"

	vramGauge := widgets.NewGauge()
	vramGauge.SetRect(0, 11, 50, 14)
	vramGauge.Title = "GPU VRAM Usage"

	gpuTempWidget := widgets.NewParagraph()
	gpuTempWidget.SetRect(0, 14, 50, 17) // Positioned below the VRAM gauge
	gpuTempWidget.Title = "GPU Temperature"

	// Set up input reader for keyboard events
	reader := bufio.NewReader(os.Stdin)

	// Polling interval (in milliseconds)
	pollingInterval := 500 * time.Millisecond
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	// Goroutine to handle updating the stats
	go func() {
		for {
			select {
			case <-ticker.C:
				// Fetch CPU clock speeds
				clockSpeeds, err := cpu.GetCPUClockSpeeds()
				if err == nil && clockSpeeds.Avg != -1 && clockSpeeds.Max != -1 {
					// Update clock speed widget to show both clock speeds and temperature
					cpuTemp, err := cpu.GetCPUTemperature()
					if err == nil && cpuTemp != "-1°C" {
						clockSpeedWidget.Text = fmt.Sprintf(
							"Average Clock: %d MHz\nMax Clock: %d MHz\nTemperature: %s",
							clockSpeeds.Avg, clockSpeeds.Max, cpuTemp,
						)
					} else {
						clockSpeedWidget.Text = fmt.Sprintf(
							"Average Clock: %d MHz\nMax Clock: %d MHz\nTemperature: N/A",
							clockSpeeds.Avg, clockSpeeds.Max,
						)
					}
				} else {
					clockSpeedWidget.Text = "CPU Clock Speed: N/A"
				}

				// Fetch CPU usage
				cpuUsage, err := usage.CalculateCPUUsage()
				if err == nil && cpuUsage.Percent != -1 {
					// Update CPU usage gauge
					usageGauge.Percent = cpuUsage.Percent
					usageGauge.Label = fmt.Sprintf("%d%%", cpuUsage.Percent)
				} else {
					usageGauge.Percent = 0
					usageGauge.Label = "CPU Usage: N/A"
				}

				// Fetch memory usage
				memStats, err := memory.GetMemoryUsage()
				if err == nil && memStats.Percent != -1 {
					// Update memory usage gauge
					memoryGauge.Percent = memStats.Percent
					memoryGauge.Label = fmt.Sprintf("%.1f GB / %.1f GB", memStats.UsedGB, memStats.TotalGB)
				} else {
					memoryGauge.Percent = 0
					memoryGauge.Label = "Memory Usage: N/A"
				}

				// Fetch GPU VRAM usage
				vramStats, err := gpu.GetVRAMUsage()
				if err == nil && vramStats.UsedMB != -1 && vramStats.TotalMB != -1 {
					// Update GPU VRAM usage gauge
					vramGauge.Percent = int(float64(vramStats.UsedMB) / float64(vramStats.TotalMB) * 100)
					vramGauge.Label = fmt.Sprintf("%d MB / %d MB", vramStats.UsedMB, vramStats.TotalMB)
				} else {
					vramGauge.Percent = 0
					vramGauge.Label = "GPU VRAM Usage: N/A"
				}

				// Fetch GPU temperature
				gpuTempData, err := temp.GetGPUTemperature()
				if err == nil && gpuTempData.EdgeTemp != "-1°C" {
					// Update GPU temperature widget
					gpuTempWidget.Text = fmt.Sprintf("Edge Temperature: %s", gpuTempData.EdgeTemp)
				} else {
					gpuTempWidget.Text = "GPU Temperature: N/A"
				}

				// Render only the widgets that have valid data
				widgetsToRender := []termui.Drawable{}
				if clockSpeedWidget.Text != "CPU Clock Speed: N/A" {
					widgetsToRender = append(widgetsToRender, clockSpeedWidget)
				}
				if usageGauge.Percent != 0 {
					widgetsToRender = append(widgetsToRender, usageGauge)
				}
				if memoryGauge.Percent != 0 {
					widgetsToRender = append(widgetsToRender, memoryGauge)
				}
				if vramGauge.Percent != 0 {
					widgetsToRender = append(widgetsToRender, vramGauge)
				}
				if gpuTempWidget.Text != "GPU Temperature: N/A" {
					widgetsToRender = append(widgetsToRender, gpuTempWidget)
				}

				// Render the widgets
				termui.Render(widgetsToRender...)
			}
		}
	}()

	// Main loop to handle keyboard input
	for {
		// Poll for keyboard input
		if input, err := reader.ReadByte(); err == nil {
			if input != 0 { // If a non-zero byte is received, exit
				fmt.Println("Exiting...")
				return
			}
		}
	}
}
