package gui

import (
	"connected/model"
	"connected/model/event"
	"connected/settings"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

var serverIsRunning = false
var clientIsRunning = false

func setupTray(desk desktop.App, w fyne.Window) {
	var startServerItem, connectClientItem *fyne.MenuItem

	startServerItem = fyne.NewMenuItem("Start Server", func() {
		if clientIsRunning {
			ShowError("You cannot start a server while connected to client!")
			return
		}
		if !serverIsRunning {
			event.GetBus().Publish(model.EventTypeServerStart)
			serverIsRunning = true
			startServerItem.Label = "Start Server"
		} else {
			event.GetBus().Publish(model.EventTypeServerStop)
			serverIsRunning = false
			startServerItem.Label = "Stop Server"
		}
		startServerItem.Checked = serverIsRunning
	})

	connectClientItem = fyne.NewMenuItem("Connect Client", func() {
		if serverIsRunning {
			ShowError("You cannot start a client while connected to server!")
			return
		}
		if clientIsRunning {
			event.GetBus().Publish(model.EventTypeClientDisconnect)
			clientIsRunning = false
			connectClientItem.Label = "Connect Client"
		} else {
			showConnectDialog()
			clientIsRunning = true
			connectClientItem.Label = "Disconnect Client"
		}
		connectClientItem.Checked = clientIsRunning
	})

	settingsItem := fyne.NewMenuItem("Settings", func() {
		showSettingsDialog()
	})

	m := fyne.NewMenu("Remote OCR",
		startServerItem,
		fyne.NewMenuItem(fmt.Sprintf("IP: %s:%d", settings.GetLocalIP(), settings.GetPort()), nil),
		connectClientItem,
		settingsItem,
	)
	desk.SetSystemTrayMenu(m)
}
