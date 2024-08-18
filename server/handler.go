package server

import (
	"bytes"
	"connected/model"
	"connected/ocr"
	"connected/security"
	"connected/settings"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net"
	"time"
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
		sendError(conn, "ERR001")
		return
	}

	// AES key for encryption
	key := security.GenerateAESKey(settings.GetPassword())

	// Continuously send OCR data to the client
	for serverRunning {
		ocrText := ocr.GetOcrResult()

		// Encrypt the OCR text
		encryptedText, err := security.Encrypt([]byte(ocrText), key)
		if err != nil {
			fmt.Println("Error encrypting data:", err)
			sendError(conn, "ERR002")
			return
		}

		err = sendData(conn, model.DataTypeOCR, encryptedText)
		if err != nil {
			fmt.Println("Error sending data to client:", err)
			return
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func sendError(conn net.Conn, code string) {
	data := model.Data{
		Type: model.DataTypeError,
		Data: code,
	}
	sendData(conn, data.Type, data.Data)
}

func sendData(conn net.Conn, dataType model.DataType, data string) error {
	response := model.Data{
		Type: dataType,
		Data: data,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		return err
	}
	jsonData = append(jsonData, '\n')
	_, err = conn.Write(jsonData)
	return err
}
