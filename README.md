# GoPerfWatch

GoPerfWatch is a CLI-based system monitoring utility written in Go for **Linux** systems. It displays real-time metrics of hardware such as CPU, GPU, and memory. 

## Features
- **CPU Monitoring**: Displays clock speed, temperature, and usage.
- **GPU Monitoring**: Tracks VRAM usage and temperature.
- **Memory Monitoring**: Shows memory usage in GB and / total on a usage bar.
- **Real-time Updates**: Data is refreshed at configurable intervals (must modify source code).

![image](https://github.com/user-attachments/assets/9e4b9658-94b4-4333-911c-673c17f50616)

## Dependencies

Before attempting install GoPerfWatch, ensure that the following packages are present on your system:

### Golang

For Ubuntu

```bash
sudo apt-get install golang-go
```

For Fedora
```bash
sudo dnf install golang
```

For Arch
```bash
sudo pacman -Syu go
```

OR download the latest Linux release from [Go's official site](https://go.dev/dl/)

### lm-sensors

For Ubuntu

```bash
sudo apt-get install lm-sensors
```

For Fedora
```bash
sudo dnf install lm_sensors
```

For Arch
```bash
sudo pacman -Syu lm_sensors
```

### glxinfo

For Ubuntu

```bash
sudo apt-get install mesa-utils
```

For Fedora
```bash
sudo dnf install glx-utils
```

For Arch
```bash
sudo pacman -Syu mesa-utils
```
## Installation

### Clone the repo

```bash
git clone https://github.com/vibraniumdroid/goperfwatch.git
cd goperfwatch
```

### Initialize the module and install termui

```bash
go mod init goperfwatch
go get github.com/gizak/termui/v3
go mod tidy
```

### Build or run directly

```bash
go build
chmod u+x goperfwatch
./goperfwatch
```
OR

```bash
go run main.go
```

## To do

* Implement GPU clock monitoring
* Verify support for existing GPU features on Nvidia GPUs
* Implement power consumption monitoring on supported hardware
* Allow for arguments when executing to enable different display modes
* Quit program only upon receiving specific input (rather than any)
* Possibly phase out some use of lm-sensors where more universal implementations may be possible
