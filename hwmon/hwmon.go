package hwmon

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/firescry/zephyr/fancurve"
	"github.com/firescry/zephyr/timeseries"
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

func listHwmon() []string {
	result, err := filepath.Glob(hwmonRootDir)
	if err != nil {
		return result
	}
	return result
}

func newDevice(hwmon string) *Device {
	device := Device{}
	device.hwmon = hwmon
	device.Name = device.readEp(hwmonNameEp)
	device.PwmMax, _ = strconv.Atoi(device.readEp(hwmonPwmMaxEp))
	device.PwmMin, _ = strconv.Atoi(device.readEp(hwmonPwmMinEp))
	device.TempSamples = timeseries.InitTimeSeries(tempSamplesNumber, device.readTemp())
	device.fanCurve = fancurve.NewCurve(50.0, 80.0, 32.0, 100.0)
	return &device
}

func (device *Device) isSupported() bool {
	return device.Name == supportedDeviceName
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
	temp, err := strconv.ParseFloat(device.readEp(hwmonTempEp), 64)
	if err != nil {
		log.Fatal(err)
	}
	return temp / 1000.0
}

// SetPwmMode sets PWM mode for device
func (device *Device) SetPwmMode(mode int) {
	device.writeEp(hwmonPwmEnableEp, mode)
	log.Printf("[%s] PWM mode change to %d\n", device.Name, mode)
}

func (device *Device) setPwm(value int) {
	device.writeEp(hwmonPwmEp, value)
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
	for _, hwmon := range listHwmon() {
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
