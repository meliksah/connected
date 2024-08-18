package gui

import (
	"connected/client"
	"connected/settings"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

func SetupAndRun() {
	a := app.New()
	w := a.NewWindow("Settings")

	// Tray Icon Setup
	if desk, ok := a.(desktop.App); ok {
		setupTray(desk, w)
	}

	w.SetCloseIntercept(func() {
		w.Hide() // Hide the window instead of closing it
	})
	w.Hide() // Start by hiding the window

	a.Run()
}

func showSettingsDialog(w fyne.Window) {
	portEntry := widget.NewEntry()
	portEntry.SetText(fmt.Sprintf("%d", settings.GetPort()))
	portEntry.SetPlaceHolder("Enter server port")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetText(settings.GetPassword())
	passwordEntry.SetPlaceHolder("Enter password")

	saveButton := widget.NewButton("Save", func() {
		port, err := strconv.Atoi(portEntry.Text)
		if err != nil {
			fmt.Println("Invalid port number. Using default port.")
			port = 8080
		}
		settings.SetPort(port)
		settings.SetPassword(passwordEntry.Text)

		err = settings.SaveSettings()
		if err != nil {
			fmt.Println("Error saving settings:", err)
		} else {
			fmt.Println("Settings saved.")
		}
		w.Hide()
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("Server Port:"),
		portEntry,
		widget.NewLabel("Password:"),
		passwordEntry,
		saveButton,
	))
	w.Show()
}

func showConnectDialog(w fyne.Window) {
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
		w.Hide()
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("Server IP:"),
		ipEntry,
		widget.NewLabel("Server Port:"),
		portEntry,
		connectButton,
	))
	w.Show()
}
