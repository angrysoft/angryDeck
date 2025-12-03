package main

import (
	"fmt"

	streamdeck "angrysoft.ovh/angry-deck/devices"
)

const VERSION = "0.1.0"

func main() {
	devs, err := streamdeck.Devices()
	if err != nil {
		fmt.Println("no Stream Deck devices found: %s", err)
	}
	if len(devs) == 0 {
		fmt.Println("no Stream Deck devices found")
	}
	for _, dev := range devs {
		fmt.Printf("Found device: %s\n", dev.ID)
	}

	d := devs[0]
	if err := d.Open(); err != nil {
		fmt.Printf("can't open device: %s\n", err)
	}

	d.Clear()
	events, err := d.ReadKeys()
	if err != nil {
		fmt.Printf("can't read keys: %s\n", err)
	}
	for e := range events {
		fmt.Printf("Key event: %+v\n", e)
	}
	defer d.Close()

	fmt.Printf("Device %s opened\n", d.ID)
}
