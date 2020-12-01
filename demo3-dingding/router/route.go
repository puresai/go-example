package router

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "local.com/sai0556/demo3-dingding/controller"
)

func Load(g *gin.Engine) *gin.Engine {
    g.Use(gin.Recovery())
    // 404
    g.NoRoute(func (c *gin.Context)  {
        c.String(http.StatusNotFound, "404 not found");
    })

    g.GET("/healthCheck", controller.HealthCheck)
	g.POST("/dingding", controller.DingDing)

    return g
}