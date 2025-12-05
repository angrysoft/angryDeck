package main

import (
	"angrysoft.ovh/angry-deck/streamdeck"
)

const VERSION = "0.1.0"

func main() {
	// devs, err := streamdeck.Devices()
	// if err != nil {
	// 	fmt.Println("no Stream Deck devices found: %s", err)
	// }
	// if len(devs) == 0 {
	// 	fmt.Println("no Stream Deck devices found")
	// }
	// for _, dev := range devs {
	// 	fmt.Printf("Found device: %s\n", dev.ID)
	// }

	// d := devs[0]
	// if err := d.Open(); err != nil {
	// 	fmt.Printf("can't open device: %s\n", err)
	// }

	// d.Clear()
	// events, err := d.ReadKeys()
	// if err != nil {
	// 	fmt.Printf("can't read keys: %s\n", err)
	// }
	// for e := range events {
	// 	fmt.Printf("Key event: %+v\n", e)
	// }
	// defer d.Close()

	// fmt.Printf("Device %s opened\n", d.ID)

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

	kesy, err := device.ListenKeys()
	if err != nil {
		println("Error listening to keys:", err.Error())
		return
	}
	for key := range kesy {
		println("Key event:", key.Index, "Pressed:", key.Pressed)
	}
}
