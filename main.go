package main

import (
	"angrysoft.ovh/angry-deck/streamdeck"
)

const VERSION = "0.1.0"

func main() {

	devices, err := streamdeck.FindDevices()
	if err != nil {
		println("Error listing devices:", err.Error())
		return
	}
	if len(devices) == 0 {
		println("No Stream Deck devices found")
		return
	}
	device := devices[0]
	println("Using device: %s - %s", device.Manufacturer, device.Product)

	err = device.OpenDeckDevice()
	if err != nil {
		println("Error opening device:", err.Error())
		return
	}
	defer device.Close()

	println("Device opened successfully")
	err = device.SetImageFromFile(0, "./test.png")
	if err != nil {
		println("Error setting image:", err.Error())
		return
	}
	println("Image set successfully")

	kesy, err := device.ListenKeys()
	if err != nil {
		println("Error listening to keys:", err.Error())
		return
	}
	for key := range kesy {
		if key.Index == 9 && !key.Pressed {
			device.Clear()
			break
		}
		println("Key event:", key.Index, "Pressed:", key.Pressed)
	}
}
