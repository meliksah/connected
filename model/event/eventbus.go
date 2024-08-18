package event

import (
	"github.com/asaskevich/EventBus"
)

// Global event bus instance
var bus = EventBus.New()

func GetBus() EventBus.Bus {
	return bus
}
