package client

import (
	"bufio"
	"bytes"
	"connected/model"
	"connected/model/event"
	"connected/security"
	"connected/settings"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
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
		Stop()
		return
	}

	hash := sha256.Sum256([]byte(settings.GetPassword()))
	encryptedMagicWord := sha256.Sum256(append(hash[:], []byte(settings.MagicWord)...))
	conn.Write(encryptedMagicWord[:])

	go func() {
		reader := bufio.NewReader(conn)
		key := security.GenerateAESKey(settings.GetPassword())

		for {
			// Read until a newline or an error occurs
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					fmt.Println("Server closed the connection")
					Stop()
					break
				}
				fmt.Println("Error reading from server:", err)
				Stop()
				break
			}

			// Trim newline characters
			line = bytes.TrimSpace(line)

			// Unmarshal the JSON response
			var response model.Data
			err = json.Unmarshal(line, &response)
			if err != nil {
				fmt.Println("Error unmarshaling server response:", err)
				continue
			}

			switch response.Type {
			case model.DataTypeError:
				event.GetBus().Publish(model.EventTypeError, response.Data)
			case model.DataTypeOCR:
				// Decrypt the OCR text
				decryptedText, err := security.Decrypt(response.Data, key)
				if err != nil {
					fmt.Println("Error decrypting data:", err)
					continue
				}
				fmt.Println("Received OCR text:", string(decryptedText)) // handle the text as needed
			}
		}
	}()
	go event.GetBus().Publish(model.EventTypeClientConnected)
}

func Stop() {
	clientRunning = false
	if conn != nil {
		conn.Close()
	}
	go event.GetBus().Publish(model.EventTypeClientDisconnected)
}

func SubscribeTopics() {
	event.GetBus().Subscribe(model.EventTypeClientConnect, Connect)
	event.GetBus().Subscribe(model.EventTypeClientDisconnect, Stop)
}

func IsRunning() bool {
	return clientRunning
}
