package server

import (
	"connected/ocr"
	"connected/settings"
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	serverRunning bool
	conn          net.Conn
	mu            sync.Mutex
)

func Start() {
	serverRunning = true

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", settings.GetPort()))
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	go func() {
		for serverRunning {
			conn, err := ln.Accept()
			if err != nil {
				if serverRunning {
					fmt.Println("Error accepting connection:", err)
				}
				break
			}
			go handleClient(conn)
		}
	}()

	go func() {
		for serverRunning {
			ocr.CaptureAndOcr(conn)
			time.Sleep(500 * time.Millisecond)
		}
	}()
}

func Stop() {
	serverRunning = false
	if conn != nil {
		conn.Close()
	}
}

func IsRunning() bool {
	return serverRunning
}
