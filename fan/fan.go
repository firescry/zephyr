package fan

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/firescry/zephyr/fancurve"
	"github.com/firescry/zephyr/hwmon"
	"github.com/firescry/zephyr/timeseries"
)

const (
	SupportedDeviceName = "amdgpu"
	TempSamplesNumber   = 30

	// PwmModeManual sets PWM mode to manual
	PwmModeManual = 1
	// PwmModeAuto sets PWM mode to auto
	PwmModeAuto = 2
)

// Device represents supported device
type Device struct {
	hwmon       string
	Name        string
	PwmMax      int
	PwmMin      int
	TempSamples *timeseries.TimeSeries
	fanCurve    *fancurve.Curve
}

func newDevice(hw string) *Device {
	device := Device{}
	device.hwmon = hw
	device.Name = device.readEp(hwmon.HwmonNameEp)
	device.PwmMax, _ = strconv.Atoi(device.readEp(hwmon.HwmonPwmMaxEp))
	device.PwmMin, _ = strconv.Atoi(device.readEp(hwmon.HwmonPwmMinEp))
	device.TempSamples = timeseries.InitTimeSeries(TempSamplesNumber, device.readTemp())
	device.fanCurve = fancurve.NewCurve(50.0, 80.0, 32.0, 100.0)
	return &device
}

func (device *Device) isSupported() bool {
	return device.Name == SupportedDeviceName
}

func (device *Device) readEp(ep string) string {
	content, err := os.ReadFile(filepath.Join(device.hwmon, ep))
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(content))
}

func (device *Device) writeEp(ep string, value int) {
	err := os.WriteFile(filepath.Join(device.hwmon, ep), []byte(fmt.Sprint(value)), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (device *Device) readTemp() float64 {
	temp, err := strconv.ParseFloat(device.readEp(hwmon.HwmonTempEp), 64)
	if err != nil {
		log.Fatal(err)
	}
	return temp / 1000.0
}

// SetPwmMode sets PWM mode for device
func (device *Device) SetPwmMode(mode int) {
	device.writeEp(hwmon.HwmonPwmEnableEp, mode)
	log.Printf("[%s] PWM mode change to %d\n", device.Name, mode)
}

func (device *Device) setPwm(value int) {
	device.writeEp(hwmon.HwmonPwmEp, value)
}

func (device *Device) percentToPwm(percent float64) (pwm int) {
	pwm = int(math.Round(float64(device.PwmMax-device.PwmMin) * percent / 100))
	switch {
	case pwm < device.PwmMin:
		return device.PwmMin
	case pwm > device.PwmMax:
		return device.PwmMax
	default:
		return pwm
	}
}

// SupportedDevices returns list of supported devices
func SupportedDevices() []*Device {
	var supportedDevices []*Device
	for _, hwmon := range hwmon.ListHwmon() {
		device := newDevice(hwmon)
		if device.isSupported() {
			log.Printf("Found supported device: %s\n", device.Name)
			supportedDevices = append(supportedDevices, device)
		} else {
			log.Printf("Found unsupported device: %s\n", device.Name)
		}
	}
	return supportedDevices
}

// Update device by reading its current temperature and adjusting PWM accordingly
func (device *Device) Update() {
	device.TempSamples.AddSample(device.readTemp())
	average := device.TempSamples.WeightedAverage()
	percent := device.fanCurve.GetPwmPercent(average)
	pwm := device.percentToPwm(percent)
	device.setPwm(pwm)
}
