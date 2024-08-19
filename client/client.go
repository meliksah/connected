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

	connection, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		go event.GetBus().Publish(model.EventTypeError, "Error connecting to server:"+err.Error())
		Stop()
		return
	}
	conn = connection

	hash := sha256.Sum256([]byte(settings.GetPassword()))
	encryptedMagicWord := sha256.Sum256(append(hash[:], []byte(settings.MagicWord)...))
	conn.Write(encryptedMagicWord[:])

	go func() {
		reader := bufio.NewReader(conn)
		key := security.GenerateAESKey(settings.GetPassword())

		for clientRunning {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					fmt.Println("Server closed the connection")
					go event.GetBus().Publish(model.EventTypeError, "Server closed the connection: "+err.Error())
					Stop()
					break
				}
				fmt.Println("Error reading from server:", err)
				go event.GetBus().Publish(model.EventTypeError, "Error reading from server: "+err.Error())
				Stop()
				break
			}

			line = bytes.TrimSpace(line)

			var response model.Data
			err = json.Unmarshal(line, &response)
			if err != nil {
				fmt.Println("Error unmarshaling server response:", err)
				continue
			}

			switch response.Type {
			case model.DataTypeError:
				go event.GetBus().Publish(model.EventTypeError, response.Data)
			case model.DataTypeOCR:
				decryptedText, err := security.Decrypt(response.Data, key)
				if err != nil {
					fmt.Println("Error decrypting data:", err)
					continue
				}
				go event.GetBus().Publish(model.EventTypeOcrTextReceived, string(decryptedText))
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
