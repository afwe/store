package mq

import (
	"deck/config"
	"github.com/streadway/amqp"
	"log"
)

var conn *amqp.Connection
var channel *amqp.Channel
var notifyClose chan *amqp.Error

func UpdateRabbitHost(host string) {
	config.RabbitURL = host
}
func Init() {
	if !config.AsyncTransferEnable {
		return
	}
	if initChannel(config.RabbitURL) {
		channel.NotifyClose(notifyClose)
	}
	go func() {
		for {
			select {
			case msg := <-notifyClose:
				conn = nil
				channel = nil
				log.Printf("rbmqclose: %+v\n", msg)
				initChannel(config.RabbitURL)
			}
		}
	}()
}
func initChannel(rabbitHost string) bool {
	if channel != nil {
		return true
	}
	conn, err := amqp.Dial(rabbitHost)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	channel, err = conn.Channel()
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
