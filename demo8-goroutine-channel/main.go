package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	c := make(chan interface{})
	v, ok := unBlockRead(c)
	fmt.Println(v, ok)

	ok = unBlockWrite(c, "joker")
	fmt.Println(ok)
}

func unBlockWrite(c chan<- interface{}, v interface{}) (ok bool) {
	select {
	case c <- v:
		return true
	case <-time.After(time.Second):
		return false
	}
}

func unBlockRead(c <-chan interface{}) (v interface{}, ok bool) {
	select {
	case v, ok = <-c:
		return
	case <-time.After(time.Second):
		return nil, false
	}
}

// func main() {
// 	timeout()
// }

// chan 配合 select
func timeout() {
	timeout := make(chan bool)
	go func() {
		time.Sleep(3e9)
		timeout <- true
	}()
	ch := make(chan int)
	select {
	case <-ch:
	case <-timeout:
		fmt.Println("timeout!")
	}
}

// 非缓冲通道，监听信号量，来自gin文档的例子
func ginHttp() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func fifo() {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	ch <- 3
	ch <- 4
	ch <- 5

	fmt.Println("1-", <-ch)
	fmt.Println("2-", <-ch)
	ch <- 6
	fmt.Println("3-", <-ch)
	fmt.Println("4-", <-ch)
	fmt.Println("5-", <-ch)
	fmt.Println("6-", <-ch)
	close(ch)
}
