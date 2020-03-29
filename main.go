package main

import (
	"github.com/firescry/zephyr/hwmon"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	supportedDevices := hwmon.SupportedDevices()

	if len(supportedDevices) == 0 {
		log.Printf("WARNING: There are no supported devices")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	loop := true

	var currentTemp float64
	var average float64

	for loop {
		select {
		case <-quit:
			loop = false
		default:
			for _, device := range supportedDevices {
				currentTemp = device.ReadTemp()
				device.TempSamples.AddSample(currentTemp)
				average = device.TempSamples.WeightedAverage()
				log.Printf("[%s] current temperature: %.1f; average: %f", device.Name, currentTemp, average)
				time.Sleep(time.Second)
			}
		}
	}
	log.Printf("Exit!")
}
