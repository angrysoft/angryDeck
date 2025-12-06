package deck

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"angrysoft.ovh/angry-deck/page"
	"angrysoft.ovh/angry-deck/streamdeck"
	"gopkg.in/yaml.v3"
)

type Deck struct {
	PagesConfigs []string `yaml:"pages_configs"`
	Default      string
	Settings     DeckSettings
	pages        map[string]*page.Page
	handlers     map[string]*page.Action
	deck         *streamdeck.DeckDevice
	configDir    string
	currentPage  string
}

type DeckSettings struct {
	Brightness int
}

const imageDir = "images"

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

	// defer device.Close()
	return &Deck{
		PagesConfigs: []string{},
		Settings:     DeckSettings{Brightness: 100},
		pages:        make(map[string]*page.Page),
		handlers:     make(map[string]*page.Action),
		deck:         &device,
	}
}

func (d *Deck) LoadDeck(path string) error {
	d.configDir = filepath.Dir(path)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, d)
	if err != nil {
		return err
	}
	for _, pageName := range d.PagesConfigs {
		page := page.NewPage()
		err := page.LoadPage(filepath.Join(d.configDir, pageName))
		if err != nil {
			return err
		}
		d.pages[page.Name] = page
		for _, button := range page.Buttons {
			d.handlers[fmt.Sprintf("%s.%d.%s", page.Name, button.Index, button.Action.OnState)] = &button.Action
		}
	}
	d.setPage(d.Default)

	return nil
}

func (d *Deck) Clear() {
	if d.deck != nil {
		d.deck.Clear()
	}
}

func (d *Deck) Close() {
	if d.deck != nil {
		d.deck.Close()
	}
}

func (d *Deck) setPage(name string) error {
	log.Println("Switching to page:", name)
	page, exists := d.pages[name]
	if !exists {
		return nil
	}
	log.Println("Set page ", name)
	d.deck.Clear()
	d.currentPage = name

	for _, button := range page.Buttons {
		fmt.Println("Setting button", button.Index, "on page", name)
		d.deck.SetButton(button.Index, d.configDir, button.Label, button.Icon)
	}
	return nil
}
