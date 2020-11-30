package mq

import (
	"github.com/streadway/amqp"
	"log"
)

var conn *amqp.Connection
var channel *amqp.Channel

var notifyClose chan *amqp.Error

func InitChannel(rabbitHost string)bool {
	if channel != nil{
		return true
	}
	conn, err := amqp.Dial(rabbitHost)
	if err != nil{
		log.Println(err.Error())
		return false
	}
	channel, err = conn.Channel()
	if err != nil{
		log.Println(err.Error())
		return false
	}
	return true
}
