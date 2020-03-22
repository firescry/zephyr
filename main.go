package main

import (
	"github.com/firescry/zephyr/hwmon"
	"github.com/firescry/zephyr/timeseries"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const timeSeriesLength = 30

// ToDo: It doesn't work :(
func waitForShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}

func handler(device *hwmon.Device, run chan bool) {
	loop := true
	var currentTemp float64
	var average float64
	weights := timeseries.GenerateWeights(timeSeriesLength)
	samples := timeseries.Initialize(timeSeriesLength, device.ReadTemp())
	for loop {
		select {
		case <-run:
			loop = false
		default:
			currentTemp = device.ReadTemp()
			samples = timeseries.AddValue(samples, currentTemp)
			average = timeseries.WeightedAverage(samples, weights)
			log.Printf("[%s] current temperature: %.1f; average: %f", device.Name, currentTemp, average)
			time.Sleep(time.Second)
		}
	}
	log.Printf("%s: FINISHED!\n", device.Name)
}

func main() {
	supportedDevices := hwmon.SupportedDevices()

	if len(supportedDevices) == 0 {
		log.Printf("WARNING: There are no supported devices")
	}

	run := make(chan bool, 1)

	for _, device := range supportedDevices {
		go handler(device, run)
	}

	waitForShutdown()

	for range supportedDevices {
		run <- false
	}
}
