package page

import (
	"log"
	"os/exec"
)

type Action struct {
	Type      string
	Value     []string
	OnRelease bool `yaml:"on_release"`
}

func (a *Action) Execute() {
	// Placeholder for action execution logic
	switch a.Type {
	case "exec":
		a.DoExec()
	default:
		log.Printf("Unknown action type: %s\n", a.Type)
	}
}

func (a *Action) DoExec() {
	log.Println("Running command")
	if len(a.Value) == 0 {
		log.Println("No command specified")
		return
	}
	var cmd *exec.Cmd
	if len(a.Value) > 1 {
		cmd = exec.Command(a.Value[0], a.Value[1:]...)
	} else {
		cmd = exec.Command(a.Value[0])

	}

	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to execute command: %v\n", err)
	}
}
