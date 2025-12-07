package main

import (
	"os"
	"os/signal"
	"syscall"

	"angrysoft.ovh/angry-deck/deck"
)

const VERSION = "0.1.0"

func main() {
	deck := deck.NewDeck()
	err := deck.LoadDeck("example_config/deck.yml")
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
