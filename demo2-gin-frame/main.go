package main

// import 这里我习惯把官方库，开源库，本地module依次分开列出
import (
	"log"
	"errors"

    "github.com/spf13/pflag"
    "github.com/spf13/viper"
    "github.com/gin-gonic/gin"

    "local.com/sai0556/demo2-gin-frame/config"
    "local.com/sai0556/demo2-gin-frame/db"
    "local.com/sai0556/demo2-gin-frame/router"
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
	
	// 连接mysql数据库
    btn := db.GetInstance().InitPool()
    if !btn {
        log.Println("init database pool failure...")
        panic(errors.New("init database pool failure"))
    }

	// redis
	db.InitRedis()

	gin.SetMode(viper.GetString("mode"))
	g := gin.New()
	g = router.Load(g)

	g.Run(viper.GetString("addr"))
}