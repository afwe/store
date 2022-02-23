package config

const (
	AsyncTransferEnable  = false
	TransExchangeName    = "upload.trans"
	TransOSSQueueName    = "upload.trans.ceph"
	TransOSSErrQueueName = "upload.trans.oss.ceph"
	TransOSSRoutingKey   = "ceph"
)

var (
	RabbitURL = "amqp://admin:admin@127.0.0.1:5672/"
)
