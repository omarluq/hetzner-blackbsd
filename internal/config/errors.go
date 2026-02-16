package config

import "fmt"

// Error is returned when configuration is invalid.
type Error struct {
	Field   string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("config error [%s]: %s", e.Field, e.Message)
}
