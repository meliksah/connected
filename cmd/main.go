package main

import (
	"connected/client"
	"connected/gui"
	"connected/ocr"
	"connected/server"
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
	subscribeTopics()

	gui.SetupAndRun()
}

func subscribeTopics() {
	gui.SubscribeTopics()
	client.SubscribeTopics()
	server.SubscribeTopics()
	ocr.SubscribeTopics()
}
