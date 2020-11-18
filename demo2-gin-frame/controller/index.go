package controller

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// 返回
type Response struct {
    Code int `json:"code"`
    Message string `json:"message"`
    Data interface{} `json:"data"`
}

// api返回结构
func ApiResponse(c *gin.Context, code int, message string, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code: code,
        Message: message,
        Data: data,
    })
}

func Index(c *gin.Context) {
    ApiResponse(c, 0, "success", nil)
}

func HealthCheck(c *gin.Context) {
    ApiResponse(c, 0, "success", nil)
}