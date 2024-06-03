package router

import (
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func SetupRouters() {
	r = gin.Default();
	r.GET("/health",func(ctx *gin.Context) { ctx.JSON(200,gin.H{"status" : "OK"}) })
}

func RunServer(listenAddr string) error {
	return r.Run(listenAddr)
}