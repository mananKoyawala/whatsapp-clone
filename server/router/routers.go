package router

import (
	"github.com/gin-gonic/gin"
	api "github.com/mananKoyawala/whatsapp-clone/internal"
	msg "github.com/mananKoyawala/whatsapp-clone/internal/message"
	user "github.com/mananKoyawala/whatsapp-clone/internal/user"
	"github.com/mananKoyawala/whatsapp-clone/internal/ws"
)

var r *gin.Engine

func SetupRouters(user *user.Handler, wshandler *ws.Handler, msgHandler *msg.Handler) {
	r = gin.Default()

	// health checking
	r.GET("/health", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"status": "OK"}) })

	// user routes
	r.POST("/users/signup", api.MakeHTTPHandleFunc(user.CreateUser))
	r.POST("/users/login", api.MakeHTTPHandleFunc(user.LoginUser))
	r.POST("/users/verify", api.MakeHTTPHandleFunc(user.VerifyUserOTP))

	// ws routes
	r.GET("/ws/connect/:uid", wshandler.WsConnector) // * for websocket connection mothod must be GET

	// msg routes
	r.POST("/msgs", api.MakeHTTPHandleFunc(msgHandler.PullAllMessages))
	r.PATCH("/msgs", api.MakeHTTPHandleFunc(msgHandler.UpdateIsReadMessage))
	r.DELETE("/msgs", api.MakeHTTPHandleFunc(msgHandler.DeleteMessage))
}

func RunServer(listenAddr string) error {
	return r.RunTLS(listenAddr, "../server.crt", "../server.key")
}
