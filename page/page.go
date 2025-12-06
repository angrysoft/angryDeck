package page

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Page struct {
	Name    string
	Buttons []Button
}

type Button struct {
	Index      uint8
	Icon       Icon
	Label      Label
	FontSize   int `yaml:"font_size"`
	Action     Action
}

type Icon struct {
	File string
	Fill string
}

type Label struct {
	Text      string
	FontSize  int    `yaml:"font_size"`
	FontColor string `yaml:"font_color"`
}

func NewPage() *Page {
	return &Page{
		Buttons: []Button{},
	}
}

func (p *Page) LoadPage(pageFile string) error {
	data, err := os.ReadFile(pageFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, p)
	if err != nil {
		return err
	}

	fmt.Printf("Loaded page: %s with %d buttons\n", p.Name, len(p.Buttons))
	return nil
}
