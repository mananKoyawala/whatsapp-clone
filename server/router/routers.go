package router

import (
	"github.com/gin-gonic/gin"
	api "github.com/mananKoyawala/whatsapp-clone/internal"
	user "github.com/mananKoyawala/whatsapp-clone/internal/user"
)

var r *gin.Engine

func SetupRouters(user *user.Handler) {
	r = gin.Default()
	r.GET("/health", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"status": "OK"}) })

	r.POST("/users/signup", api.MakeHTTPHandleFunc(user.CreateUser))
	r.POST("/users/login", api.MakeHTTPHandleFunc(user.LoginUser))
	r.POST("/users/verify", api.MakeHTTPHandleFunc(user.VerifyUserOTP))
}

func RunServer(listenAddr string) error {
	return r.Run(listenAddr)
}
