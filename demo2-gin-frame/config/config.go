package config

import (
	"fmt"

    "github.com/fsnotify/fsnotify"
    "github.com/spf13/viper"
)

type Config struct {
    Name string
}

// 对外的初始化配置方法
func Run(cfg string) error {
    c := Config{
        Name: cfg,
    }

    if err := c.init(); err != nil {
        return err
    }

    c.watchConfig()

    return nil
}

func (c *Config) init() error {
    if c.Name != "" {
        viper.SetConfigFile(c.Name)
    } else {
        // 默认配置文件是./config.yaml
        viper.AddConfigPath(".")
        viper.SetConfigName("config")
    }

    viper.SetConfigType("yaml")
    // viper解析配置文件
    err := viper.ReadInConfig() 
    if err != nil {
        panic(fmt.Errorf("Fatal error config file: %s \n", err))
    }

	// 简单打印下配置
    fmt.Println(viper.GetString("name"))

    return nil
}

func (c *Config) watchConfig() {
    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
        fmt.Println("Config file changed:", e.Name)
    })
}
