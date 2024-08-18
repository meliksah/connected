package client

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"net"
	"remote_ocr/settings"
)

var (
	clientRunning bool
	conn          net.Conn
)

func Connect(w fyne.Window) {
	clientRunning = true

	ip, port := settings.GetLastIP(), settings.GetLastPort()
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}

	hash := sha256.Sum256([]byte(settings.GetPassword()))
	encryptedMagicWord := sha256.Sum256(append(hash[:], []byte(settings.MagicWord)...))
	conn.Write(encryptedMagicWord[:])

	go func() {
		scanner := bufio.NewScanner(conn)
		ocrText := widget.NewMultiLineEntry()
		for scanner.Scan() {
			text := scanner.Text()
			ocrText.SetText(text)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading from server:", err)
		}
	}()
}

func Stop() {
	clientRunning = false
	if conn != nil {
		conn.Close()
	}
}

func IsRunning() bool {
	return clientRunning
}
