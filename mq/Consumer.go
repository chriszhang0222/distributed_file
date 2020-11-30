package mq

import "log"

var done chan bool

func StartConsume(qName, cName string, callback func(msg []byte)bool){
	msgs, err := channel.Consume(qName, cName, true, false, false, false, nil)
	if err != nil{
		log.Fatal(err)
		return
	}

	done = make(chan bool)
	go func() {
		for d := range msgs{
			processErr := callback(d.Body)
			if processErr {
				log.Println(processErr)
			}
		}
	}()
	<-done
	channel.Close()
}

func StopConsume(){
	done <- true
}
