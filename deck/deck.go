package deck

import (
	"log"
	"os"

	"angrysoft.ovh/angry-deck/streamdeck"
	"gopkg.in/yaml.v3"
)

type Deck struct {
	Pages        []string
	DefaultPage  string
	Settings     DeckSettings
	pageHandlers map[string]*Page
	deck         *streamdeck.DeckDevice
}

type DeckSettings struct {
	Brightness int
}

func NewDeck() *Deck {
	devices, err := streamdeck.FindDevices()
	if err != nil {
		log.Println("Error listing devices:", err.Error())
		return nil
	}
	if len(devices) == 0 {
		log.Println("No Stream Deck devices found")
		return nil
	}
	device := devices[0]
	log.Println("Using device:", device.Manufacturer, device.Product)
	err = device.OpenDeckDevice()
	if err != nil {
		log.Println("Error opening device:", err.Error())
		return nil
	}

	defer device.Close()
	return &Deck{
		Pages:        []string{},
		Settings:     DeckSettings{Brightness: 100},
		pageHandlers: make(map[string]*Page),
		deck:         &device,
	}
}

func (d *Deck) LoadDeck(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, d)
	if err != nil {
		return err
	}
	for _, pageName := range d.Pages {
		page := newPage()
		err := page.LoadPage(pageName)
		if err != nil {
			return err
		}
		d.pageHandlers[pageName] = page
	}

	return nil
}

func (d *Deck) setPage(name string) error {
	page, exists := d.pageHandlers[name]
	if !exists {
		return nil
	}
	for _, button := range page.Buttons {
		if button.Image != "" {
			err := d.deck.SetImageFromFile(button.Index, button.Image)
			if err != nil {
				log.Println("Error setting image for button", button.Index, ":", err.Error())
			}
		}
		if button.Label != "" {
			err := d.deck.SetImageWithTextFromFile(button.Index, button.Image, button.Label)
			if err != nil {
				log.Println("Error setting label for button", button.Index, ":", err.Error())
			}
		}
	}
	return nil
}
