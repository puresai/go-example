package main

// import 这里我习惯把官方库，开源库，本地module依次分开列出
import (
	"fmt"
	"time"
	"strconv"
	"net/http"

    "github.com/spf13/pflag"
    "github.com/spf13/viper"
    "github.com/gin-gonic/gin"

    "local.com/sai0556/demo2-gin-frame/config"
    "local.com/sai0556/demo2-gin-frame/db"
    "local.com/sai0556/demo2-gin-frame/router"
    "local.com/sai0556/demo2-gin-frame/logger"
    "local.com/sai0556/demo2-gin-frame/graceful"
)

var (
    conf = pflag.StringP("config", "c", "", "config filepath")
)

func main() {
    pflag.Parse()

    // 初始化配置
    if err := config.Run(*conf); err != nil {
        panic(err)
	}

	logger.Info("i'm log123-----Info")
	logger.Error("i'm log123-----Error")

	
	// 连接mysql数据库
	DB := db.GetDB()
	defer db.CloseDB(DB)

	// redis
	db.InitRedis()

	go func() {
		pingServer()
	}()

	gin.SetMode(viper.GetString("mode"))
	g := gin.New()
	g = router.Load(g)

	// g.Run(viper.GetString("addr"))


	// logger.Info("启动http服务端口%s\n", viper.GetString("addr"))

	// time.Sleep(2*time.Second)
	if err := graceful.ListenAndServe(viper.GetString("addr"), g); err != nil && err != http.ErrServerClosed {
		logger.Error("fail:http服务启动失败: %s\n", err)
	}
}

// 健康检查
// func pingServer() error {
// 	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
// 		url := fmt.Sprintf("%s%s%s", "http://127.0.0.1", viper.GetString("addr"), viper.GetString("healthCheck"))
// 		fmt.Println(url)
// 		resp, err := http.Get(url)
// 		if err == nil && resp.StatusCode == 200 {
// 			return nil
// 		}
// 		time.Sleep(time.Second)
// 	}
// 	return errors.New("健康检测404")
// }

// 健康检查
func pingServer() {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		url := fmt.Sprintf("%s%s%s", "http://127.0.0.1", viper.GetString("addr"), viper.GetString("healthCheck"))
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			fmt.Println("health check success!")
			return
		}
		fmt.Println("check fail -" + strconv.Itoa(i+1)+"times")
		time.Sleep(time.Second)
	}
	fmt.Println("Cannot connect to the router!!!")
}