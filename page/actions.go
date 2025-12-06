package page

import (
	"log"
	"os/exec"
)

type Action struct {
	Type    string
	Value   string
	OnState string
}

func (a *Action) Execute() {
	// Placeholder for action execution logic
	switch a.Type {
	case "exec":
		a.doExec()
	default:
		log.Printf("Unknown action type: %s\n", a.Type)
	}
}

func (a *Action) doExec() {
	cmd := exec.Command(a.Value)
	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to execute command: %v\n", err)
	}
}
