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

	// Load error messages
	err = settings.LoadErrors()
	if err != nil {
		fmt.Println("Error loading error messages:", err)
		return
	}

	gui.SetupAndRun()
}
