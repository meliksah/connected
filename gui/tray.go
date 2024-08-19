package gui

import (
	"connected/model"
	"connected/model/event"
	"connected/settings"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

var (
	serverIsRunning     = false
	clientIsRunning     = false
	startServerItem     *fyne.MenuItem
	connectClientItem   *fyne.MenuItem
	connectedClientItem *fyne.MenuItem
	runningServerItem   *fyne.MenuItem
	clientSubMenu       *fyne.Menu
	trayMenu            *fyne.Menu
	desktopApp          desktop.App
)

func setupTray(desk desktop.App, w fyne.Window) {
	desktopApp = desk // Store the desktop app reference for later use

	// Initialize the menu items
	startServerItem = fyne.NewMenuItem("Start Server", toggleServer)
	connectClientItem = fyne.NewMenuItem("Connect Client", toggleClient)

	settingsItem := fyne.NewMenuItem("Settings", func() {
		showSettingsDialog()
	})

	trayMenu = fyne.NewMenu("Remote OCR",
		startServerItem,
		connectClientItem,
		settingsItem,
	)

	desk.SetSystemTrayMenu(trayMenu)

	// Subscribe to events to update menu items
	subscribeToEvents()
}

func toggleServer() {
	if clientIsRunning {
		ShowError("You cannot start a server while connected to client!")
		return
	}
	if !serverIsRunning {
		event.GetBus().Publish(model.EventTypeServerStart)
	} else {
		event.GetBus().Publish(model.EventTypeServerStop)
	}
}

func toggleClient() {
	if serverIsRunning {
		ShowError("You cannot start a client while connected to server!")
		return
	}
	if !clientIsRunning {
		showConnectDialog()
	} else {
		event.GetBus().Publish(model.EventTypeClientDisconnect)
	}
}

func subscribeToEvents() {
	event.GetBus().SubscribeAsync(model.EventTypeServerStarted, func() {
		serverIsRunning = true
		startServerItem.Label = "Stop Server"

		// Add running server item directly after the Start Server item
		if runningServerItem == nil {
			runningServerItem = fyne.NewMenuItem(fmt.Sprintf("Server running at: %s:%d", settings.GetLocalIP(), settings.GetPort()), nil)
			clientSubMenu = fyne.NewMenu("Connected Clients")
			runningServerItem.ChildMenu = clientSubMenu
			insertMenuItemAfter(startServerItem, runningServerItem)
		}
	}, false)

	event.GetBus().SubscribeAsync(model.EventTypeServerStopped, func() {
		serverIsRunning = false
		startServerItem.Label = "Start Server"

		// Remove running server item
		removeMenuItem(runningServerItem)
		runningServerItem = nil
		clientSubMenu = nil
	}, false)

	event.GetBus().SubscribeAsync(model.EventTypeClientConnected, func() {
		clientIsRunning = true
		connectClientItem.Label = "Disconnect Client"

		// Add connected client item directly after the Connect Client item
		if connectedClientItem == nil {
			connectedClientItem = fyne.NewMenuItem(fmt.Sprintf("Connected to: %s:%d", settings.GetLastIP(), settings.GetLastPort()), func() {
				if isOcrWindowHidden == true {
					ocrWindow.Show()
				} else {
					ocrWindow.Hide()
				}
			})
			insertMenuItemAfter(connectClientItem, connectedClientItem)
		}
	}, false)

	event.GetBus().SubscribeAsync(model.EventTypeClientDisconnected, func() {
		clientIsRunning = false
		connectClientItem.Label = "Connect Client"

		// Remove connected client item
		removeMenuItem(connectedClientItem)
		connectedClientItem = nil
	}, false)

	event.GetBus().SubscribeAsync(model.EventTypeServerClientListUpdated, updateClientList, false)
}

func updateClientList(clientlist []string) {
	if clientSubMenu == nil {
		return
	}

	clientSubMenu.Items = nil // Clear existing items

	for _, ip := range clientlist {
		clientItem := fyne.NewMenuItem(ip, func() {
			event.GetBus().Publish(model.EventTypeServerDisconnectClient, ip)
		})
		clientSubMenu.Items = append(clientSubMenu.Items, clientItem)
	}

	desktopApp.SetSystemTrayMenu(trayMenu) // Update the system tray menu
}

func insertMenuItemAfter(afterItem *fyne.MenuItem, newItem *fyne.MenuItem) {
	for i, menuItem := range trayMenu.Items {
		if menuItem == afterItem {
			trayMenu.Items = append(trayMenu.Items[:i+1], append([]*fyne.MenuItem{newItem}, trayMenu.Items[i+1:]...)...)
			desktopApp.SetSystemTrayMenu(trayMenu) // Update the system tray menu
			break
		}
	}
}

func removeMenuItem(item *fyne.MenuItem) {
	if item == nil {
		return
	}
	for i, menuItem := range trayMenu.Items {
		if menuItem == item {
			trayMenu.Items = append(trayMenu.Items[:i], trayMenu.Items[i+1:]...)
			desktopApp.SetSystemTrayMenu(trayMenu) // Update the system tray menu
			break
		}
	}
}
