package model

const (
	EventTypeError                   string = "error:occurred"
	EventTypeServerStart             string = "server:start"
	EventTypeServerStarted           string = "server:started"
	EventTypeServerStop              string = "server:stop"
	EventTypeServerStopped           string = "server:stopped"
	EventTypeClientConnect           string = "client:connect"
	EventTypeClientConnected         string = "client:connected"
	EventTypeClientDisconnect        string = "client:disconnect"
	EventTypeClientDisconnected      string = "client:disconnected"
	EventTypeServerClientListUpdated string = "server:client:updated"
	EventTypeServerDisconnectClient  string = "server:client:disconnect"
)
