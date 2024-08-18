package gui

import (
	"connected/settings"
	"fmt"
	"fyne.io/fyne/v2/dialog"
)

func HandleError(errorCode string) {
	dialog.ShowError(fmt.Errorf(settings.GetErrorMessage(errorCode)), GetGuiWindow())
}
