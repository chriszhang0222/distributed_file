package mq

import (
	"distributed_file/config"
	"github.com/streadway/amqp"
	"log"
)

func Publish(exchange, routingKey string, msg []byte)bool{
	if !InitChannel(config.RabbitURL){
		return false
	}
	if nil == channel.Publish(exchange, routingKey, false, false, amqp.Publishing{ContentType: "text/plain",
		Body: msg}){
		return true
	}
	log.Println("error")
	return false
}
