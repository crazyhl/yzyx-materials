package bus

import "github.com/asaskevich/EventBus"

var Bus EventBus.Bus

func Init() {
	Bus = EventBus.New()
}
