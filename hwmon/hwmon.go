package hwmon

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

const supportedDeviceName = "amdgpu"
const hwmonRootDir = "/sys/class/hwmon"
const hwmonNameEp = "name"
const hwmonTempEp = "temp1_input"

type Device struct {
	Name      string
	DeviceDir string
	epName    string
	epTemp    string
}

func joinPaths(p, q string) string {
	p = strings.TrimRight(p, "/")
	q = strings.TrimLeft(q, "/")
	return strings.Join([]string{p, q}, "/")
}

func readFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(content))
}

func writeFile(path string, value string) {
	err := ioutil.WriteFile(path, []byte(value), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func listDir(path string) []string {
	var r []string
	items, err := ioutil.ReadDir(path)
	if err != nil {
		return r
	}
	for _, item := range items {
		r = append(r, item.Name())
	}
	return r
}

func (device *Device) ReadTemp() (temp float64) {
	temp, _ = strconv.ParseFloat(readFile(device.epTemp), 64)
	temp = temp / 1000.0
	return temp
}

func createDevice(hwmonDir string) *Device {
	device := Device{}
	device.DeviceDir = hwmonDir
	device.epName = joinPaths(hwmonDir, hwmonNameEp)
	device.epTemp = joinPaths(hwmonDir, hwmonTempEp)
	return &device
}

func (device *Device) setDeviceName() {
	device.Name = readFile(device.epName)
}

func (device *Device) isSupported() bool {
	return readFile(device.epName) == supportedDeviceName
}

func SupportedDevices() []*Device {
	var supportedDevices []*Device
	allHwmonDevices := listDir(hwmonRootDir)
	for _, hwmonDevice := range allHwmonDevices {
		hwmonDeviceDir := joinPaths(hwmonRootDir, hwmonDevice)
		device := createDevice(hwmonDeviceDir)
		device.setDeviceName()
		if device.isSupported() {
			log.Printf("Found supported device: %s [%s]\n", device.Name, device.DeviceDir)
			supportedDevices = append(supportedDevices, device)
		} else {
			log.Printf("Found unsupported device: %s [%s]\n", device.Name, device.DeviceDir)
		}
	}
	return supportedDevices
}
