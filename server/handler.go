package server

import (
	"bytes"
	"connected/ocr"
	"connected/settings"
	"crypto/sha256"
	"fmt"
	"net"
)

func handleClient(conn net.Conn) {
	defer conn.Close()

	clientMagicWord := make([]byte, sha256.Size)
	_, err := conn.Read(clientMagicWord)
	if err != nil {
		fmt.Println("Error reading magic word from client:", err)
		return
	}

	hash := sha256.Sum256([]byte(settings.GetPassword()))
	expectedMagicWord := sha256.Sum256(append(hash[:], []byte(settings.MagicWord)...))

	if !bytes.Equal(clientMagicWord, expectedMagicWord[:]) {
		fmt.Println("Password mismatch. Connection refused.")
		conn.Write([]byte("Password mismatch.\n"))
		return
	}

	mu.Lock()
	ocrText := ocr.GetOcrResult()
	mu.Unlock()

	_, err = conn.Write([]byte(ocrText + "\n"))
	if err != nil {
		fmt.Println("Error sending data to client:", err)
		return
	}
}
