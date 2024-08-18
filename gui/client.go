package gui

import (
	"connected/client"
	"connected/settings"
	"fmt"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

func showConnectDialog() {
	ipEntry := widget.NewEntry()
	ipEntry.SetText(settings.GetLastIP())
	ipEntry.SetPlaceHolder("Enter server IP")

	portEntry := widget.NewEntry()
	portEntry.SetText(fmt.Sprintf("%d", settings.GetLastPort()))
	portEntry.SetPlaceHolder("Enter server port")

	connectButton := widget.NewButton("Connect", func() {
		ip := ipEntry.Text
		port, err := strconv.Atoi(portEntry.Text)
		if err != nil {
			fmt.Println("Invalid port number.")
			return
		}

		// Save the latest IP and port to settings
		settings.SetLastIP(ip)
		settings.SetLastPort(port)
		err = settings.SaveSettings()
		if err != nil {
			fmt.Println("Error saving settings:", err)
			return
		}

		// Connect to the server
		client.Connect(ip, port)
		gui_window.Hide()
	})

	gui_window.SetContent(container.NewVBox(
		widget.NewLabel("Server IP:"),
		ipEntry,
		widget.NewLabel("Server Port:"),
		portEntry,
		connectButton,
	))
	gui_window.Show()
}
