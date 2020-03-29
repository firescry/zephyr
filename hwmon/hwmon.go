package hwmon

import (
	"fmt"
	"github.com/firescry/zephyr/timeseries"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	tempSamplesNumber   = 30
	supportedDeviceName = "amdgpu"
	hwmonRootDir        = "/sys/class/drm/card[0-9]/device/hwmon/hwmon[0-9]"
	hwmonNameEp         = "name"
	hwmonPwmEp          = "pwm1"
	hwmonPwmEnableEp    = "pwm1_enable"
	hwmonPwmMaxEp       = "pwm1_max"
	hwmonPwmMinEp       = "pwm1_min"
	hwmonTempEp         = "temp1_input"
)

type Device struct {
	hwmon       string
	Name        string
	PwmMax      string
	PwmMin      string
	TempSamples *timeseries.TimeSeries
}

func listHwmon() []string {
	result, err := filepath.Glob(hwmonRootDir)
	if err != nil {
		return result
	}
	return result
}

func createDevice(hwmon string) *Device {
	device := Device{}
	device.hwmon = hwmon
	device.Name = device.readEp(hwmonNameEp)
	device.PwmMax = device.readEp(hwmonPwmMaxEp)
	device.PwmMin = device.readEp(hwmonPwmMinEp)
	device.TempSamples = timeseries.InitTimeSeries(tempSamplesNumber, device.ReadTemp())
	return &device
}

func (device *Device) isSupported() bool {
	return device.Name == supportedDeviceName
}

func (device *Device) readEp(ep string) string {
	content, err := ioutil.ReadFile(filepath.Join(device.hwmon, ep))
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(content))
}

func (device *Device) writeEp(ep string, value int) {
	err := ioutil.WriteFile(ep, []byte(fmt.Sprint(value)), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (device *Device) ReadTemp() float64 {
	temp, err := strconv.ParseFloat(device.readEp(hwmonTempEp), 64)
	if err != nil {
		log.Fatal(err)
	}
	return temp / 1000.0
}

func (device *Device) SetMode(mode int) {
	device.writeEp(hwmonPwmEnableEp, mode)
}

func (device *Device) SetPwm(value int) {
	device.writeEp(hwmonPwmEp, value)
}

func SupportedDevices() []*Device {
	var supportedDevices []*Device
	for _, hwmon := range listHwmon() {
		device := createDevice(hwmon)
		if device.isSupported() {
			log.Printf("Found supported device: %s\n", device.Name)
			supportedDevices = append(supportedDevices, device)
		} else {
			log.Printf("Found unsupported device: %s\n", device.Name)
		}
	}
	return supportedDevices
}
