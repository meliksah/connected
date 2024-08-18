package ocr

import (
	"bytes"
	"fmt"
	"image/png"
	"net"
	"sync"

	"github.com/kbinani/screenshot"
	"github.com/otiai10/gosseract/v2"
)

var mu sync.Mutex
var ocrResult string

func CaptureAndOcr(conn net.Conn) {
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		fmt.Println("Error capturing screen:", err)
		return
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		fmt.Println("Error encoding image to PNG:", err)
		return
	}

	client := gosseract.NewClient()
	defer client.Close()
	client.SetImageFromBytes(buf.Bytes())
	text, err := client.Text()
	if err != nil {
		fmt.Println("Error performing OCR:", err)
		return
	}

	mu.Lock()
	ocrResult = text
	mu.Unlock()

	if conn != nil {
		_, err = conn.Write([]byte(text + "\n"))
		if err != nil {
			fmt.Println("Error sending OCR text to client:", err)
		}
	}
}

func GetOcrResult() string {
	return ocrResult
}
