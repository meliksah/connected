package gui

import (
	"connected/client"
	"connected/server"
	"connected/settings"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func setupTray(desk desktop.App, w fyne.Window) {
	var startServerItem, connectClientItem *fyne.MenuItem

	startServerItem = fyne.NewMenuItem("Start Server", func() {
		if server.IsRunning() {
			server.Stop()
			startServerItem.Label = "Start Server"
		} else {
			server.Start()
			startServerItem.Label = "Stop Server"
		}
		startServerItem.Checked = server.IsRunning()
	})

	connectClientItem = fyne.NewMenuItem("Connect Client", func() {
		if client.IsRunning() {
			client.Stop()
			connectClientItem.Label = "Connect Client"
		} else {
			showConnectDialog()
			connectClientItem.Label = "Disconnect Client"
		}
		connectClientItem.Checked = client.IsRunning()
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
