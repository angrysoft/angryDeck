package main

import (
	"os"
	"os/signal"
	"syscall"

	"angrysoft.ovh/angry-deck/deck"
)

const VERSION = "0.1.0"

var defaultConfig = "example_config/deck.yml"

func main() {
	if len(os.Args) > 1 {
		defaultConfig = os.Args[1]
	}
	deck := deck.NewDeck()
	err := deck.LoadDeck(defaultConfig)
	if err != nil {
		println("Error loading deck:", err.Error())
		return
	}
	defer func() {
		deck.Clear()
		deck.Close()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		println()
		println("Received signal:", sig.String())
		deck.Clear()
		deck.Close()
		os.Exit(1)
	}()

	deck.ListHandlers()
	deck.Listen()

	println("Exiting Angry Deck")
}
