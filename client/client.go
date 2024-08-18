package client

import (
	"bufio"
	"connected/settings"
	"crypto/sha256"
	"fmt"
	"net"
)

var (
	clientRunning bool
	conn          net.Conn
)

func Connect(ip string, port int) {
	clientRunning = true

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
		for scanner.Scan() {
			text := scanner.Text()
			fmt.Println("Received OCR text:", text) // handle the text as needed
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
