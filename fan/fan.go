package fan

import (
	"log"
	"math"

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
	hwmon       hwmon.Hwmon
	PwmMax      int
	PwmMin      int
	TempSamples *timeseries.TimeSeries
	fanCurve    *fancurve.Curve
}

func newDevice(hw hwmon.Hwmon) *Device {
	return &Device{
		hwmon:       hw,
		PwmMax:      hw.GetPwmMax(),
		PwmMin:      hw.GetPwmMin(),
		TempSamples: timeseries.InitTimeSeries(TempSamplesNumber, hw.GetTemp()),
		fanCurve:    fancurve.NewCurve(50.0, 80.0, 32.0, 100.0),
	}
}

func (device *Device) isSupported() bool {
	return device.hwmon.GetName() == SupportedDeviceName
}

// SetPwmMode sets PWM mode for device
func (device *Device) SetPwmMode(mode int) {
	device.hwmon.SetPwmMode(mode)
	log.Printf("[%s] PWM mode change to %d\n", device.hwmon.GetName(), mode)
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
	for _, hwmon := range hwmon.FindAll() {
		device := newDevice(hwmon)
		if device.isSupported() {
			log.Printf("Found supported device: %s\n", device.hwmon.GetName())
			supportedDevices = append(supportedDevices, device)
		} else {
			log.Printf("Found unsupported device: %s\n", device.hwmon.GetName())
		}
	}
	return supportedDevices
}

// Update device by reading its current temperature and adjusting PWM accordingly
func (device *Device) Update() {
	device.TempSamples.AddSample(device.hwmon.GetTemp())
	average := device.TempSamples.WeightedAverage()
	percent := device.fanCurve.GetPwmPercent(average)
	pwm := device.percentToPwm(percent)
	device.hwmon.SetPwm(pwm)
}
