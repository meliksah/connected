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
	ln            net.Listener
	mu            sync.Mutex
	clients       = make(map[string]net.Conn) // Track clients by IP address
)

func Start() {
	serverRunning = true

	var err error
	ln, err = net.Listen("tcp", fmt.Sprintf(":%d", settings.GetPort()))
	if err != nil {
		fmt.Println("Error starting server:", err)
		Stop() // Properly stop the server if there's an error
		return
	}

	localAddr, _ := GetLocalIP()
	settings.SetLocalIP(localAddr)

	go func() {
		for serverRunning {
			conn, err := ln.Accept()
			if err != nil {
				if serverRunning {
					fmt.Println("Error accepting connection:", err)
				}
				continue
			}
			clientIP := conn.RemoteAddr().(*net.TCPAddr).IP.String()
			mu.Lock()
			clients[clientIP] = conn
			mu.Unlock()
			updateClientList()
			go handleClient(conn, clientIP) // Handle each client in a separate goroutine
		}
	}()

	go event.GetBus().Publish(model.EventTypeServerStarted)
}

func Stop() {
	mu.Lock()
	defer mu.Unlock()

	// Disconnect all clients
	for _, conn := range clients {
		conn.Close()
	}
	clients = make(map[string]net.Conn) // Clear the client list

	if ln != nil {
		err := ln.Close() // Properly close the listener
		if err != nil {
			fmt.Println("Error closing listener:", err)
		}
		ln = nil // Reset the listener
	}

	serverRunning = false
	go event.GetBus().Publish(model.EventTypeServerStopped)
}

func DisconnectClient(ip string) {
	mu.Lock()
	defer mu.Unlock()

	if conn, exists := clients[ip]; exists {
		conn.Close()
		delete(clients, ip)
		updateClientList()
	}
}

func GetClients() []string {
	ips := make([]string, 0, len(clients))
	for ip := range clients {
		ips = append(ips, ip)
	}
	return ips
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
	event.GetBus().Subscribe(model.EventTypeServerDisconnectClient, DisconnectClient)
}

func IsRunning() bool {
	return serverRunning
}

func GetLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

func updateClientList() {
	go event.GetBus().Publish(model.EventTypeServerClientListUpdated, GetClients())
}
