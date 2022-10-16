package hwmon

import (
	"path/filepath"
)

const (
	HwmonRootDir     = "/sys/class/drm/card[0-9]/device/hwmon/hwmon[0-9]"
	HwmonNameEp      = "name"
	HwmonPwmEp       = "pwm1"
	HwmonPwmEnableEp = "pwm1_enable"
	HwmonPwmMaxEp    = "pwm1_max"
	HwmonPwmMinEp    = "pwm1_min"
	HwmonTempEp      = "temp1_input"
)

func ListHwmon() []string {
	result, err := filepath.Glob(HwmonRootDir)
	if err != nil {
		return result
	}
	return result
}
