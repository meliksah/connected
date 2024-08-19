package settings

// ErrorSettings contains a map of error codes to their corresponding messages.
type ErrorSettings struct {
	Errors map[string]string
}

// Initialize the error settings with predefined error messages.
var errorSettings = ErrorSettings{
	Errors: map[string]string{
		"ERR001": "Password mismatch. Connection refused.",
		"ERR002": "Invalid OCR data received.",
		"ERR003": "Failed to process OCR.",
	},
}

// GetErrorMessage returns the error message corresponding to the given error code.
// If the error code does not exist, it returns an empty string.
func GetErrorMessage(code string) string {
	if msg, exists := errorSettings.Errors[code]; exists {
		return msg
	}
	return ""
}
