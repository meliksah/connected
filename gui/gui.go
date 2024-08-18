package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"remote_ocr/settings"
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
