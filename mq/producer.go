package mq

import (
	"deck/config"
	"github.com/streadway/amqp"
	"log"
)

func Publish(exchange, routineKey string, msg []byte) bool {
	if !initChannel(config.RabbitURL) {
		return false
	}
	err := channel.Publish(
		exchange,
		routineKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "test/plain",
			Body:        msg,
		})
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
