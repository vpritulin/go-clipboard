package main

import (
	"clipboard/manager"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"golang.design/x/clipboard"
	"log"
)

func main() {
	err := clipboard.Init()
	if err != nil {
		log.Fatalf("Failed to initialize clipboard: %v", err)
	}

	cm := manager.NewClipboardManager()
	newEntries := make(chan manager.ClipboardEntry)

	go cm.MonitorClipboard(newEntries)

	a := app.New()
	mainWindow := a.NewWindow("Clipboard Manager")

	icon, err := fyne.LoadResourceFromPath("icon.png")
	if err != nil {
		log.Printf("Failed to load tray icon: %v", err)
	} else {
		a.SetIcon(icon)
	}

	ui, historyList := cm.CreateUI(mainWindow)

	go func() {
		for range newEntries {
			historyList.Refresh()
		}
	}()

	mainWindow.SetContent(ui)
	mainWindow.Resize(fyne.NewSize(475, 500))
	settingsManager := manager.NewSettings(a)

	if desk, ok := a.(desktop.App); ok {
		trayMenu := fyne.NewMenu("Clipboard Manager",
			fyne.NewMenuItem("Show History", func() {
				mainWindow.Show()
			}),
			fyne.NewMenuItem("Settings", func() {
				settingsManager.ShowSettingsWindow()
			}),
			fyne.NewMenuItem("Quit", func() {
				a.Quit()
			}),
		)
		desk.SetSystemTrayMenu(trayMenu)
	}

	mainWindow.SetCloseIntercept(func() {
		mainWindow.Hide()
	})

	a.Run()
}
