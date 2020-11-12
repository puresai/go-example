package main

import (
	"time"
	"net/http"
	"encoding/json"
)


// 自定义返回
type JsonRes struct {
	Code int `json:"code"`
	Data interface{} `json:"data"`
	Msg string `json:"msg"`
	TimeStamp int64 `json:"timestmap"`
}

func apiResult(w http.ResponseWriter, code int, data interface{}, msg string) {
	body, _ := json.Marshal(JsonRes{
		Code: code, 
		Data: data, 
		Msg: msg, 
		// 获取时间戳
		TimeStamp: time.Now().Unix(),
	})
    w.Write(body)
}



func main()  {
	srv := http.Server{
		Addr: ":8080",
		Handler: http.TimeoutHandler(http.HandlerFunc(defaultHttp), 2*time.Second, "Timeout!!!"),
	}
	srv.ListenAndServe()
}

// 默认http处理
func defaultHttp(w http.ResponseWriter, r *http.Request) {
	path, httpMethod := r.URL.Path, r.Method

	if path == "/" {
		w.Write([]byte("index"))
		return 
	}

	if path == "/hello" && httpMethod == "POST" {
		sayHello(w, r)
		return 
	}

	if path == "/sleep" {
		// 模拟一下业务处理超时
		time.Sleep(4*time.Second)
		return 
	}

	if path == "/path" {
		w.Write([]byte("path:"+path+", method:"+httpMethod))
		return 
	}

	// 自定义404
	http.Error(w, "you lost???", http.StatusNotFound)
}

// 处理hello，并接收参数输出json
func sayHello(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

    // 第一种方式，但是没有name参数会报错
    // name := query["name"][0]

    // 第二种方式
	name := query.Get("name")
	
    apiResult(w, 0, name+" say "+r.PostFormValue("some"), "success")
}
