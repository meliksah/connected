package server

import (
	"connected/model"
	"connected/model/event"
	"connected/settings"
	"fmt"
	"net"
	"sync"
)

var (
	serverRunning bool
	mu            sync.Mutex
)

func Start() {
	serverRunning = true

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", settings.GetPort()))
	if err != nil {
		fmt.Println("Error starting server:", err)
		Stop()
		return
	}

	go func() {
		for serverRunning {
			conn, err := ln.Accept()
			if err != nil {
				if serverRunning {
					fmt.Println("Error accepting connection:", err)
				}
				continue
			}
			go handleClient(conn) // Handle each client in a separate goroutine
		}
	}()

	go event.GetBus().Publish(model.EventTypeServerStarted)
}

func Stop() {
	serverRunning = false
	go event.GetBus().Publish(model.EventTypeServerStopped)
}

func SubscribeTopics() {
	err := event.GetBus().Subscribe(model.EventTypeServerStart, Start)
	if err != nil {
		return
	}
	err2 := event.GetBus().Subscribe(model.EventTypeServerStop, Stop)
	if err2 != nil {
		return
	}
}

func IsRunning() bool {
	return serverRunning
}
