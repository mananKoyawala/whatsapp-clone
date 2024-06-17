package router

import (
	"github.com/gin-gonic/gin"
	api "github.com/mananKoyawala/whatsapp-clone/internal"
	"github.com/mananKoyawala/whatsapp-clone/internal/contact"
	msg "github.com/mananKoyawala/whatsapp-clone/internal/message"
	user "github.com/mananKoyawala/whatsapp-clone/internal/user"
	"github.com/mananKoyawala/whatsapp-clone/internal/ws"
	"github.com/mananKoyawala/whatsapp-clone/middleware"
	"github.com/mananKoyawala/whatsapp-clone/service/upload"
)

var r *gin.Engine

func SetupRouters(user *user.Handler, wshandler *ws.Handler, msgHandler *msg.Handler, uploadHandler *upload.AwsHandler, contact *contact.Handler) {
	r = gin.Default()

	// health checking
	r.GET("/health", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"status": "OK"}) })

	// user routes
	r.POST("/users/signup", api.MakeHTTPHandleFunc(user.CreateUser))
	r.POST("/users/login", api.MakeHTTPHandleFunc(user.LoginUser))
	r.POST("/users/verify", api.MakeHTTPHandleFunc(user.VerifyUserOTP))

	// middleware
	r.Use(middleware.AuthMiddleware(user))

	// ws routes
	r.GET("/ws/connect/:uid", wshandler.WsConnector) // * for websocket connection mothod must be GET

	// msg routes
	r.POST("/msgs", api.MakeHTTPHandleFunc(msgHandler.PullAllMessages))
	r.PATCH("/msgs", api.MakeHTTPHandleFunc(msgHandler.UpdateIsReadMessage))
	r.DELETE("/msgs", api.MakeHTTPHandleFunc(msgHandler.DeleteMessage))

	// file routes
	r.POST("/upload", api.MakeHTTPHandleFunc(uploadHandler.UploaFile))
	r.DELETE("/delete", api.MakeHTTPHandleFunc(uploadHandler.DeleteFile))

	// contact routes
	r.POST("/contacts", api.MakeHTTPHandleFunc(contact.AddContact))
	r.GET("/contacts/:id", api.MakeHTTPHandleFunc(contact.GetContacts))

}

func RunServer(listenAddr string) error {
	return r.RunTLS(listenAddr, "../server.crt", "../server.key")
}
