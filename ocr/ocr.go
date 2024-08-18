package ocr

import (
	"bytes"
	"connected/model"
	"connected/model/event"
	"fmt"
	"image/png"
	"sync"
	"time"

	"github.com/kbinani/screenshot"
	"github.com/otiai10/gosseract/v2"
)

var (
	mu        sync.Mutex
	ocrResult string
	running   bool
	stopChan  chan struct{}
)

// CaptureAndOcr captures the screen, performs OCR, and updates the result.
func CaptureAndOcr() {
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
}

// GetOcrResult returns the latest OCR result.
func GetOcrResult() string {
	mu.Lock()
	defer mu.Unlock()
	return ocrResult
}

// startOcrProcessing starts a goroutine that continuously captures OCR data every 500 milliseconds.
func StartOcrProcessing() {
	if running {
		return
	}
	running = true
	stopChan = make(chan struct{})

	go func() {
		for {
			select {
			case <-stopChan:
				return
			default:
				CaptureAndOcr()
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
}

// stopOcrProcessing stops the OCR capturing process.
func StopOcrProcessing() {
	if !running {
		return
	}
	close(stopChan)
	running = false
}

// SubscribeTopics subscribes to server start/stop events.
func SubscribeTopics() {
	err := event.GetBus().Subscribe(model.EventTypeServerStarted, StartOcrProcessing)
	if err != nil {
		return
	}
	err2 := event.GetBus().Subscribe(model.EventTypeServerStopped, StopOcrProcessing)
	if err2 != nil {
		return
	}
}
