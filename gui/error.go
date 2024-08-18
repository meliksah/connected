package gui

import (
	"connected/settings"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func HandleError(errorCode string) {
	ShowError(settings.GetErrorMessage(errorCode))
}

func ShowError(errorMessage string) {
	// Create a new window for the error message
	errorWindow := GetGuiApp().NewWindow("Error")

	// Create a label for the error message
	errorLabel := widget.NewLabel(errorMessage)

	// Create a button to close the error window
	closeButton := widget.NewButton("Close", func() {
		errorWindow.Close()
	})

	// Create a vertical box container for the error message and button
	content := container.NewVBox(
		errorLabel,
		closeButton,
	)

	// Set the content of the window
	errorWindow.SetContent(content)

	// Set the window size to a small size since it only contains the error message and a button
	errorWindow.Resize(fyne.NewSize(300, 100))

	// Show the window as a standalone error window
	errorWindow.Show()
}
