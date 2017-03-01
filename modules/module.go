package modules

import "time"

// Module serves as the base for all modules
type Module interface {
	Info(output chan string, trigger chan string)
	Initialize()
	IsTriggered(test string) bool

	GetAlignment() string
	GetInterval() time.Duration
}
