package hwmon

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	hwmonRootDir = "/sys/class/hwmon/*"

	epName    = "name"
	epPwm     = "pwm1"
	epPwmMax  = "pwm1_max"
	epPwmMin  = "pwm1_min"
	epPwmMode = "pwm1_enable"
	epTemp    = "temp1_input"
)

type Hwmon struct {
	epNamePath    string
	epPwmPath     string
	epPwmMaxPath  string
	epPwmMinPath  string
	epPwmModePath string
	epTempPath    string
}

func NewHwmon(path string) Hwmon {
	return Hwmon{
		epNamePath:    filepath.Join(path, epName),
		epPwmPath:     filepath.Join(path, epPwm),
		epPwmMaxPath:  filepath.Join(path, epPwmMax),
		epPwmMinPath:  filepath.Join(path, epPwmMin),
		epPwmModePath: filepath.Join(path, epPwmMode),
		epTempPath:    filepath.Join(path, epTemp),
	}
}

func FindAll() []Hwmon {
	paths, _ := filepath.Glob(hwmonRootDir)
	results := make([]Hwmon, len(paths))
	for _, path := range paths {
		result := NewHwmon(path)
		results = append(results, result)
	}
	return results
}

func (hwmon Hwmon) GetName() string {
	data, _ := os.ReadFile(hwmon.epNamePath)
	return strings.TrimSpace(string(data))
}

func (hwmon Hwmon) GetPwmMax() int {
	data, _ := os.ReadFile(hwmon.epPwmMaxPath)
	result, _ := strconv.Atoi(string(data))
	return result
}

func (hwmon Hwmon) GetPwmMin() int {
	data, _ := os.ReadFile(hwmon.epPwmMinPath)
	result, _ := strconv.Atoi(string(data))
	return result
}

func (hwmon Hwmon) GetTemp() float64 {
	data, _ := os.ReadFile(hwmon.epTempPath)
	result, _ := strconv.ParseFloat(string(data), 64)
	return result / 1000
}

func (hwmon Hwmon) SetPwm(pwm int) {
	data := []byte(strconv.Itoa(pwm))
	os.WriteFile(hwmon.epPwmPath, data, 0644)
}

func (hwmon Hwmon) SetPwmMode(mode int) {
	data := []byte(strconv.Itoa(mode))
	os.WriteFile(hwmon.epPwmPath, data, 0644)
}
