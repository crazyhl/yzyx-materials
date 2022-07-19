package account

import "github.com/crazyhl/yzyx-materials/internal/bus"

func InitBus() {
	bus.Bus.Subscribe("account:updateProfit", UpdateAccountProfit)
}

func DestoryBus() {
	bus.Bus.Unsubscribe("account:updateProfit", UpdateAccountProfit)
}
