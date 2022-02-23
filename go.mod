module deck

go 1.14

require (
	github.com/garyburd/redigo v1.6.3
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/contrib v0.0.0-20201101042839-6a891bf89f19
	github.com/gin-gonic/gin v1.7.7
	github.com/go-redis/redis/v8 v8.11.4
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/json-iterator/go v1.1.12
	github.com/juju/ratelimit v1.0.1
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins/registry/consul v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/go-plugins/registry/kubernetes v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/go-plugins/wrapper/breaker/hystrix v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/go-plugins/wrapper/monitoring/prometheus v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2 v2.9.1 // indirect
	github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/micro v1.18.0 // indirect
	github.com/mitchellh/mapstructure v1.4.3
	github.com/moxiaomomo/go-bindata-assetfs v1.0.0
	github.com/prometheus/client_golang v1.12.1
	github.com/streadway/amqp v1.0.0
	google.golang.org/grpc/examples v0.0.0-20220210231334-75fd0240ac41 // indirect
	google.golang.org/protobuf v1.27.1
	gopkg.in/amz.v1 v1.0.0-20150111123259-ad23e96a31d2
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
