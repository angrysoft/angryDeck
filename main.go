package main

import (
	"angrysoft.ovh/angry-deck/deck"
)

const VERSION = "0.1.0"

func main() {
	deck := deck.NewDeck()
	err := deck.LoadDeck("deck.yml")
	if err != nil {
		println("Error loading deck:", err.Error())
		return
	}

	// println("Device opened successfully")
	// err = device.SetImageFromFile(0, "./test.png")
	// if err != nil {
	// 	println("Error setting image:", err.Error())
	// 	return
	// }
	// err = device.SetImageWithTextFromFile(1, "./test.png", "Hello World")
	// if err != nil {
	// 	println("Error setting image:", err.Error())
	// 	return
	// }
	// println("Image set successfully")

	// kesy, err := device.ListenKeys()
	// if err != nil {
	// 	println("Error listening to keys:", err.Error())
	// 	return
	// }
	// for key := range kesy {
	// 	if key.Index == 9 && !key.Pressed {
	// 		device.Clear()
	// 		break
	// 	}
	// 	println("Key event:", key.Index, "Pressed:", key.Pressed)
	// }
}
