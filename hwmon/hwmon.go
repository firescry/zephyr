package hwmon

import (
	"fmt"
	"github.com/firescry/zephyr/fancurve"
	"github.com/firescry/zephyr/timeseries"
	"io/ioutil"
	"log"
	"math"
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
	PwmModeManual       = 1
	PwmModeAuto         = 2
)

type device struct {
	hwmon       string
	Name        string
	PwmMax      int
	PwmMin      int
	TempSamples *timeseries.TimeSeries
	fanCurve    func(float64) float64
}

func listHwmon() []string {
	result, err := filepath.Glob(hwmonRootDir)
	if err != nil {
		return result
	}
	return result
}

func newDevice(hwmon string) *device {
	device := device{}
	device.hwmon = hwmon
	device.Name = device.readEp(hwmonNameEp)
	device.PwmMax, _ = strconv.Atoi(device.readEp(hwmonPwmMaxEp))
	device.PwmMin, _ = strconv.Atoi(device.readEp(hwmonPwmMinEp))
	device.TempSamples = timeseries.InitTimeSeries(tempSamplesNumber, device.readTemp())
	device.fanCurve = fancurve.Curve
	return &device
}

func (device *device) isSupported() bool {
	return device.Name == supportedDeviceName
}

func (device *device) readEp(ep string) string {
	content, err := ioutil.ReadFile(filepath.Join(device.hwmon, ep))
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(content))
}

func (device *device) writeEp(ep string, value int) {
	err := ioutil.WriteFile(filepath.Join(device.hwmon, ep), []byte(fmt.Sprint(value)), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (device *device) readTemp() float64 {
	temp, err := strconv.ParseFloat(device.readEp(hwmonTempEp), 64)
	if err != nil {
		log.Fatal(err)
	}
	return temp / 1000.0
}

func (device *device) SetPwmMode(mode int) {
	device.writeEp(hwmonPwmEnableEp, mode)
	log.Printf("[%s] PWM mode change to %d\n", device.Name, mode)
}

func (device *device) setPwm(value int) {
	device.writeEp(hwmonPwmEp, value)
}

func (device *device) percentToPwm(percent float64) (pwm int) {
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

func SupportedDevices() []*device {
	var supportedDevices []*device
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

func (device *device) Update() {
	device.TempSamples.AddSample(device.readTemp())
	average := device.TempSamples.WeightedAverage()
	percent := device.fanCurve(average)
	pwm := device.percentToPwm(percent)
	device.setPwm(pwm)
}
