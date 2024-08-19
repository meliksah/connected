package gui

import (
	"connected/model"
	"connected/model/event"
	"connected/settings"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

var gui_app fyne.App
var gui_window fyne.Window

func SetupAndRun() {
	gui_app = app.New()
	gui_window = gui_app.NewWindow("Settings")

	// Tray Icon Setup
	if desk, ok := gui_app.(desktop.App); ok {
		setupTray(desk, gui_window)
	}

	gui_window.SetCloseIntercept(func() {
		gui_window.Hide() // Hide the window instead of closing it
	})
	gui_window.Hide() // Start by hiding the window

	gui_app.Run()
}

func SubscribeTopics() {
	event.GetBus().SubscribeAsync(model.EventTypeError, HandleError, false)
	event.GetBus().SubscribeAsync(model.EventTypeOcrTextReceived, HandleOcrTextReceivedEvent, false)
}

func showSettingsDialog() {
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
		gui_window.Hide()
	})

	gui_window.SetContent(container.NewVBox(
		widget.NewLabel("Server Port:"),
		portEntry,
		widget.NewLabel("Password:"),
		passwordEntry,
		saveButton,
	))
	gui_window.Show()
}

func GetGuiApp() fyne.App {
	return gui_app
}

func GetGuiWindow() fyne.Window {
	return gui_window
}
