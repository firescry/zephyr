package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/firescry/zephyr/fan"
)

func main() {
	supportedDevices := fan.SupportedDevices()

	if len(supportedDevices) == 0 {
		log.Printf("WARNING: There are no supported devices")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for _, device := range supportedDevices {
		device.SetPwmMode(fan.PwmModeManual)
	}

	loop := true
	for loop {
		select {
		case <-quit:
			loop = false
		default:
			for _, device := range supportedDevices {
				device.Update()
				time.Sleep(time.Second)
			}
		}
	}

	for _, device := range supportedDevices {
		device.SetPwmMode(fan.PwmModeAuto)
	}
}
