package greeting

import "fmt"

// GetWelcomeMessage processes the core business logic for greetings
func GetWelcomeMessage(name string) string {
	if name == "" {
		return "Hello World!"
	}
	return fmt.Sprintf("Hello %s!", name)
}
