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

	d.deck.SetBrightness(uint8(d.Settings.Brightness))

	for _, pageName := range d.PagesConfigs {
		page := page.NewPage()
		err := page.LoadPage(filepath.Join(d.configDir, pageName))
		if err != nil {
			return err
		}
		d.pages[page.Name] = page
		for _, button := range page.Buttons {
			onState := "pressed"
			if button.Action.OnRelease {
				onState = "released"
			}
			d.handlers[fmt.Sprintf("%s.%d.%s", page.Name, button.Index, onState)] = &button.Action
		}
	}
	d.setPage(d.Default)

	return nil
}

func (d *Deck) SetBrightness(brightness uint8) {
	if d.deck != nil {
		d.deck.SetBrightness(brightness)
	}
}

func (d *Deck) ListHandlers() {
	for key, action := range d.handlers {
		log.Println("Handler:", key, "Action Type:", action.Type, "Value:", action.Value)
	}
}

func (d *Deck) Listen() {
	if d.deck != nil {
		event, err := d.deck.ListenKeys()
		if err != nil {
			println("Error listening to keys:", err.Error())
			return
		}
		for ev := range event {
			println("Key event:", ev.Index, "Pressed:", ev.Pressed)
			state := "released"
			if ev.Pressed {
				state = "pressed"
			}
			actionTrigger := fmt.Sprintf("%s.%d.%s", d.currentPage, ev.Index, state)

			action, exists := d.getAction(actionTrigger)
			log.Println("Action:", actionTrigger, exists)
			if exists {
				if action.Type == "set_page" {
					err := d.setPage(action.Value[0])
					if err != nil {
						log.Println("Error setting page:", err.Error())
					}
				} else {
					action.DoExec()
				}
			}
		}
	}
}

func (d *Deck) Clear() {
	if d.deck != nil {
		log.Println("Clear Deck")
		d.deck.Clear()
	}
}

func (d *Deck) Close() {
	if d.deck != nil {
		log.Println("Close Deck")
		d.deck.Close()
	}
}

func (d *Deck) getAction(key string) (*page.Action, bool) {
	action, exists := d.handlers[key]
	return action, exists
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
