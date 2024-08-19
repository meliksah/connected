package gui

import (
	"connected/model"
	"connected/model/event"
	"connected/settings"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

var ocrWindow fyne.Window
var connectionWindow fyne.Window
var ocrText *widget.Entry
var isFocused bool
var refreshButton *widget.Button
var isOcrWindowHidden bool

func HandleOcrTextReceivedEvent(decryptedText string) {
	if ocrWindow == nil {
		ocrWindow = GetGuiApp().NewWindow("OCR Text")
		ocrText = widget.NewMultiLineEntry()
		ocrText.SetText(decryptedText)
		ocrText.OnChanged = func(text string) {
			if !isFocused {
				ocrText.SetText(text)
			}
		}
		refreshButton = widget.NewButton("Pause Refresh", func() {
			isFocused = !isFocused
			if isFocused {
				refreshButton.SetText("Pause Refresh")
			} else {
				refreshButton.SetText("Resume Refresh")
			}
		})
		closeButton := widget.NewButton("Close", func() {
			isOcrWindowHidden = true
			ocrWindow.Hide()
		})

		buttons := container.NewHBox(
			refreshButton,
			closeButton,
		)

		buttonsContainer := container.NewVBox(
			buttons,
		)

		buttonsContainer.Resize(fyne.NewSize(500, 50)) // Set a fixed height for the button row

		ocrWindow.SetContent(container.NewBorder(
			nil,              // No top widget
			buttonsContainer, // Bottom is the buttons container
			nil,              // No left widget
			nil,              // No right widget
			ocrText,          // Center is the OCR text
		))

		ocrWindow.Resize(fyne.NewSize(500, 500))
		ocrWindow.SetCloseIntercept(func() {
			isOcrWindowHidden = true
			ocrWindow.Hide()
		})
		isOcrWindowHidden = false
		ocrWindow.Show()
	} else {
		if !isFocused {
			ocrText.SetText(decryptedText)
		}
	}
}

func showConnectDialog() {
	if connectionWindow == nil {
		connectionWindow = GetGuiApp().NewWindow("Connect to Client")
		connectionWindow.SetCloseIntercept(func() {
			connectionWindow.Hide()
		})
	}
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

		event.GetBus().Publish(model.EventTypeClientConnect, ip, port)
		connectionWindow.Hide()
	})

	connectionWindow.SetContent(container.NewVBox(
		widget.NewLabel("Server IP:"),
		ipEntry,
		widget.NewLabel("Server Port:"),
		portEntry,
		connectButton,
	))
	connectionWindow.Resize(fyne.NewSize(250, 150))
	connectionWindow.CenterOnScreen()
	connectionWindow.Show()
}
