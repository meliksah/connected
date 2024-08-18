package client

import (
	"bufio"
	"connected/model"
	"connected/model/event"
	"connected/settings"
	"crypto/sha256"
	"encoding/json"
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
			var response model.Data
			err := json.Unmarshal(scanner.Bytes(), &response)
			if err != nil {
				fmt.Println("Error unmarshaling server response:", err)
				continue
			}

			switch response.Type {
			case model.DataTypeError:
				event.GetBus().Publish(model.EventTypeError, response.Data)
			case model.DataTypeOCR:
				fmt.Println("Received OCR text:", response.Data) // handle the text as needed
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading from server:", err)
		}
	}()
	event.GetBus().Publish(model.EventTypeClientConnected)

}

func Stop() {
	clientRunning = false
	if conn != nil {
		conn.Close()
	}
	event.GetBus().Publish(model.EventTypeClientDisconnected)
}

func SubscribeTopics() {
	event.GetBus().Subscribe(model.EventTypeClientConnect, Connect)
	event.GetBus().Subscribe(model.EventTypeClientDisconnect, Stop)
}

func IsRunning() bool {
	return clientRunning
}
