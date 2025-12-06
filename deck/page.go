package deck

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
	Background string
	Image      string
	Label      string
	Action     Action
}

func newPage() *Page {
	return &Page{
		Buttons: []Button{},
	}
}

func (p *Page) LoadPage(path string) error {
	data, err := os.ReadFile(path)
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
