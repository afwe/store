package mq

import "log"

var done chan bool

func StartConsumer(qName, cName string, callback func(msg []byte) bool) {
	msgs, err := channel.Consume(
		qName,
		cName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Print(err.Error())
		return
	}
	done := make(chan bool)
	go func() {
		for msg := range msgs {
			processErr := callback(msg.Body)
			if processErr {
				done <- true
			}
		}
	}()
	<-done
	channel.Close()
}
func StopConsume() {
	done <- true
}
