package main

import (
	"connected/gui"
	"connected/settings"
	"fmt"
)

func main() {
	// Load settings
	err := settings.LoadSettings()
	if err != nil {
		fmt.Println("Error loading settings:", err)
		return
	}

	gui.SetupAndRun()
}
