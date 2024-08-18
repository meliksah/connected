package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"image/png"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/kbinani/screenshot"
	"github.com/otiai10/gosseract/v2"
)

type Settings struct {
	ServerPort int    `json:"server_port"`
	ServerIP   string `json:"server_ip"`
	Password   string `json:"password"`
	LastIP     string `json:"last_ip"`
	LastPort   int    `json:"last_port"`
}

type AppState struct {
	serverRunning bool
	clientRunning bool
	settings      Settings
	ocrResult     string
	mu            sync.Mutex
	conn          net.Conn
	settingsPath  string
}

var state = AppState{
	serverRunning: false,
	clientRunning: false,
	settingsPath:  filepath.Join(os.Getenv("HOME"), ".remote_ocr", "settings.json"),
}

const MagicWord = "PASSWORD_TEXT"

func main() {
	// Load settings
	err := loadSettings()
	if err != nil {
		fmt.Println("Error loading settings:", err)
		return
	}

	a := app.New()
	w := a.NewWindow("Settings")
	state.settingsPath = getLocalIP()

	// Tray Icon Setup
	if desk, ok := a.(desktop.App); ok {

		var startServerItem, connectClientItem *fyne.MenuItem

		startServerItem = fyne.NewMenuItem("Start Server", func() {
			if state.serverRunning {
				stopServer()
				startServerItem.Label = "Start Server"
			} else {
				startServer()
				startServerItem.Label = "Stop Server"
			}
			startServerItem.Checked = state.serverRunning
		})

		connectClientItem = fyne.NewMenuItem("Connect Client", func() {
			if state.clientRunning {
				stopClient()
				connectClientItem.Label = "Connect Client"
			} else {
				connectClient()
				connectClientItem.Label = "Disconnect Client"
			}
			connectClientItem.Checked = state.clientRunning
		})

		settingsItem := fyne.NewMenuItem("Settings", func() {
			showSettingsDialog(w)
		})

		m := fyne.NewMenu("Remote OCR",
			startServerItem,
			fyne.NewMenuItem(fmt.Sprintf("IP: %s:%d", state.settingsPath, state.settings.ServerPort), nil),
			connectClientItem,
			settingsItem,
		)
		desk.SetSystemTrayMenu(m)
	}

	w.SetCloseIntercept(func() {
		w.Hide() // Hide the window instead of closing it
	})
	w.Hide() // Start by hiding the window

	a.Run()
}

func startServer() {
	state.serverRunning = true

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", state.settings.ServerPort))
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	go func() {
		for state.serverRunning {
			conn, err := ln.Accept()
			if err != nil {
				// If the listener is closed, break the loop
				if state.serverRunning {
					fmt.Println("Error accepting connection:", err)
				}
				break
			}
			state.conn = conn
			go handleClient(conn)
		}
	}()

	go func() {
		for state.serverRunning {
			captureAndOcr()
			time.Sleep(500 * time.Millisecond)
		}
	}()
}

func stopServer() {
	state.serverRunning = false
	if state.conn != nil {
		state.conn.Close()
	}
}

func connectClient() {
	state.clientRunning = true
	w := fyne.CurrentApp().NewWindow("OCR Client")

	ipEntry := widget.NewEntry()
	ipEntry.SetText(state.settings.LastIP)
	ipEntry.SetPlaceHolder("Enter server IP")

	portEntry := widget.NewEntry()
	portEntry.SetText(fmt.Sprintf("%d", state.settings.LastPort))
	portEntry.SetPlaceHolder("Enter server port")

	connectButton := widget.NewButton("Connect", func() {
		serverIP := ipEntry.Text
		serverPort, err := strconv.Atoi(portEntry.Text)
		if err != nil {
			fmt.Println("Invalid port number.")
			return
		}

		// Save the last connected IP and port to settings
		state.settings.LastIP = serverIP
		state.settings.LastPort = serverPort
		err = saveSettings()
		if err != nil {
			fmt.Println("Error saving settings:", err)
		}

		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			return
		}
		state.conn = conn

		// Send encrypted magic word to the server for verification
		hash := sha256.Sum256([]byte(state.settings.Password))
		encryptedMagicWord := sha256.Sum256(append(hash[:], []byte(MagicWord)...))
		conn.Write(encryptedMagicWord[:])

		go func() {
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				text := scanner.Text()
				// Update the OCR text in the main UI thread
				ocrText := widget.NewMultiLineEntry()
				ocrText.SetText(text)
			}
			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading from server:", err)
			}
		}()

		copyButton := widget.NewButton("Copy to Clipboard", func() {
			w.Clipboard().SetContent(ipEntry.Text)
		})

		w.SetContent(container.NewVBox(
			widget.NewLabel("OCR Text from Server:"),
			ocrText,
			copyButton,
		))

		w.Resize(fyne.NewSize(400, 300))
		w.Show()
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("Server IP:"),
		ipEntry,
		widget.NewLabel("Server Port:"),
		portEntry,
		connectButton,
	))
	w.Resize(fyne.NewSize(400, 200))
	w.Show()
}

func stopClient() {
	state.clientRunning = false
	if state.conn != nil {
		state.conn.Close()
	}
}

func captureAndOcr() {
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		fmt.Println("Error capturing screen:", err)
		return
	}

	// Convert the image to PNG format
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

	state.mu.Lock()
	state.ocrResult = text
	state.mu.Unlock()

	if state.conn != nil {
		_, err = state.conn.Write([]byte(text + "\n"))
		if err != nil {
			fmt.Println("Error sending OCR text to client:", err)
		}
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// Read the encrypted magic word from the client
	clientMagicWord := make([]byte, sha256.Size)
	_, err := conn.Read(clientMagicWord)
	if err != nil {
		fmt.Println("Error reading magic word from client:", err)
		return
	}

	// Encrypt the magic word with the server's password
	hash := sha256.Sum256([]byte(state.settings.Password))
	expectedMagicWord := sha256.Sum256(append(hash[:], []byte(MagicWord)...))

	// Compare the received magic word with the expected magic word
	if !bytes.Equal(clientMagicWord, expectedMagicWord[:]) {
		fmt.Println("Password mismatch. Connection refused.")
		conn.Write([]byte("Password mismatch.\n"))
		return
	}

	for state.serverRunning {
		state.mu.Lock()
		ocrText := state.ocrResult
		state.mu.Unlock()

		_, err := conn.Write([]byte(ocrText + "\n"))
		if err != nil {
			fmt.Println("Error sending data to client:", err)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func showSettingsDialog(w fyne.Window) {
	portEntry := widget.NewEntry()
	portEntry.SetText(fmt.Sprintf("%d", state.settings.ServerPort))
	portEntry.SetPlaceHolder("Enter server port")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetText(state.settings.Password)
	passwordEntry.SetPlaceHolder("Enter password")

	saveButton := widget.NewButton("Save", func() {
		state.settings.ServerPort = func() int {
			port, err := strconv.Atoi(portEntry.Text)
			if err != nil {
				fmt.Println("Invalid port number. Using default port.")
				return 8080
			}
			return port
		}()

		state.settings.Password = passwordEntry.Text
		err := saveSettings()
		if err != nil {
			fmt.Println("Error saving settings:", err)
		} else {
			fmt.Println("Settings saved.")
		}
		w.Hide()
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("Server Port:"),
		portEntry,
		widget.NewLabel("Password:"),
		passwordEntry,
		saveButton,
	))
	w.Show()
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func loadSettings() error {
	file, err := os.Open(state.settingsPath)
	if os.IsNotExist(err) {
		state.settings = Settings{ServerPort: 8080, Password: "", LastIP: "127.0.0.1", LastPort: 8080}
		err = saveSettings()
		if err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&state.settings)
	if err != nil {
		return err
	}

	return nil
}

func saveSettings() error {
	err := os.MkdirAll(filepath.Dir(state.settingsPath), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(state.settingsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&state.settings)
	if err != nil {
		return err
	}

	return nil
}
