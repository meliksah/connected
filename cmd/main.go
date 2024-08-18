package main

import (
	"fmt"
	"remote_ocr/gui"
	"remote_ocr/settings"
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
