package main

import (
	"gospacex-tz/srv/basic/RabbitMQ"
	_ "gospacex-tz/srv/basic/inits"
)

func main() {
	RabbitMQ.ConsumeStockDeduct()
}
