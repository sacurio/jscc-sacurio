package util

const (
	DefaultAppName string = "ChatBot"
	DefaultPort    string = "8085"
)

// ValidateStringNotEmpty ensures a variable value is not empty, return a default value provided if it is.
func ValidateStringNotEmpty(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}

	return value
}
