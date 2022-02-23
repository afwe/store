package route

import (
	"deck/assets"
	"deck/service/apigw/handler"
	"deck/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/static"
	assetfs "github.com/moxiaomomo/go-bindata-assetfs"

	"github.com/gin-gonic/gin"
)

var (
	customMonitor = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "custom_monitor",
		Help: "custom monitor some data",
	}, []string{"msg", "time", "program"})
)

//初始化Prometheus模型
func init() {
	prometheus.MustRegister(customMonitor)
}

type binaryFileSystem struct {
	fs http.FileSystem
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{
		Asset:     assets.Asset,
		AssetDir:  assets.AssetDir,
		AssetInfo: assets.AssetInfo,
		Prefix:    root,
	}
	return &binaryFileSystem{
		fs,
	}
}

// Router : 网关api路由
func Router() *gin.Engine {
	router := gin.Default()

	//	router.Static("/static/", "./static")
	// 将静态文件打包到bin文件
	router.Use(static.Serve("/static/", BinaryFileSystem("static")))
	promMonitor := util.NewPrometheusMonitor("user", "apigw")
	router.Use(promMonitor.PromMiddleware())
	// 注册
	router.GET("/user/signup", handler.SignupHandler)
	router.POST("/user/signup", handler.DoSignupHandler)
	// 登录
	router.GET("/user/signin", handler.SigninHandler)
	router.POST("/user/signin", handler.DoSigninHandler)

	router.GET("/metrics", prometheusHandler())
	router.Use(handler.Authorize())

	// 用户查询
	router.POST("/user/info", handler.UserInfoHandler)

	// 用户文件查询
	router.POST("/file/query", handler.FileQueryHandler)
	// 用户文件修改(重命名)
	router.POST("/file/update", handler.FileMetaUpdateHandler)

	return router
}
