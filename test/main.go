package main

import (
	"distributed_file/mq"
	"fmt"
)
import "distributed_file/config"
func main(){

	mq.InitChannel(config.RabbitURL)
	mq.StartConsume(config.TransOSSQueueName, config.TransExchangeName, func(msg []byte) bool {
		fmt.Println(string(msg))
		return true

	})

}

