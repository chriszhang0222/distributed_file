package config

const (
	AsyncTransferEnable = false
	TransExchangeName = "uploadserver.trans"
	TransOSSQueueName = "uploadserver.trans.oss"
	TransOSSErrQueueName = "uploadserver.trans.oss.err"
	TransOSSRoutingKey = "oss"

)

var (
	RabbitURL = "amqp://admin:admin@192.168.0.10/my_vhost"
)
