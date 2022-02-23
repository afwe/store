package main

import (
	"deck/common"
	"deck/service/db/config"
	dbConn "deck/service/db/conn"
	dbProxy "deck/service/db/proto"
	dbRpc "deck/service/db/rpc"
	"deck/util"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func startRpcService() {
	service := micro.NewService(
		micro.Name("go.micro.service.dbproxy"), // 在注册中心中的服务名称
		micro.Flags(common.CustomFlags...),
		micro.Registry(util.ConsulRegistry()),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)

	service.Init(
		micro.Action(func(c *cli.Context) {
			// 检查是否指定dbhost
			dbhost := c.String("dbhost")
			if len(dbhost) > 0 {
				log.Println("custom db address: " + dbhost)
				config.UploadDBHost(dbhost)
			}
		}),
	)

	// 初始化db connection
	dbConn.InitDBConn()

	dbProxy.RegisterDBServiceHandler(service.Server(), new(dbRpc.DB))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
	go startPrometheus()
}

func main() {
	startRpcService()
}
func startPrometheus() {
	http.Handle("/metrics", promhttp.Handler())
	// 启动web服务，监听8085端口
	go func() {
		err := http.ListenAndServe("localhost:8881", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
}

// res, err := mapper.FuncCall("/user/UserExist", []interface{}{"haha"}...)
// log.Printf("error: %+v\n", err)
// log.Printf("result: %+v\n", res[0].Interface())

// res, err = mapper.FuncCall("/user/UserExist", []interface{}{"admin"}...)
// log.Printf("error: %+v\n", err)
// log.Printf("result: %+v\n", res[0].Interface())
