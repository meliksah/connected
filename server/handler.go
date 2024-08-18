package server

import (
	"bytes"
	"connected/model"
	"connected/ocr"
	"connected/settings"
	"crypto/sha256"
	"encoding/json"
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
		sendError(conn, "ERR001")
		return
	}

	mu.Lock()
	ocrText := ocr.GetOcrResult()
	mu.Unlock()

	sendData(conn, model.DataTypeOCR, ocrText)
}

func sendError(conn net.Conn, code string) {
	data := model.Data{
		Type: model.DataTypeError,
		Data: code,
	}
	sendData(conn, data.Type, data.Data)
}

func sendData(conn net.Conn, dataType model.DataType, data string) {
	response := model.Data{
		Type: dataType,
		Data: data,
	}
	jsonData, _ := json.Marshal(response)
	conn.Write(jsonData)
}
