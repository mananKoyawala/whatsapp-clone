package router

import (
	"github.com/gin-gonic/gin"
	api "github.com/mananKoyawala/whatsapp-clone/internal"
	"github.com/mananKoyawala/whatsapp-clone/internal/contact"
	"github.com/mananKoyawala/whatsapp-clone/internal/group"
	msg "github.com/mananKoyawala/whatsapp-clone/internal/message"
	user "github.com/mananKoyawala/whatsapp-clone/internal/user"
	"github.com/mananKoyawala/whatsapp-clone/internal/ws"
	"github.com/mananKoyawala/whatsapp-clone/middleware"
	"github.com/mananKoyawala/whatsapp-clone/service/upload"
)

var r *gin.Engine

func SetupRouters(user *user.Handler, wshandler *ws.Handler, msgHandler *msg.Handler, uploadHandler *upload.AwsHandler, contact *contact.Handler, groupHand *group.Handler) {
	r = gin.Default()

	// health checking
	r.GET("/health", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"status": "OK"}) })

	// user routes
	r.POST("/users/signup", api.MakeHTTPHandleFunc(user.CreateUser))
	r.POST("/users/login", api.MakeHTTPHandleFunc(user.LoginUser))
	r.POST("/users/verify", api.MakeHTTPHandleFunc(user.VerifyUserOTP))
	r.POST("/users/token", api.MakeHTTPHandleFunc(user.RefreshToken))

	// middleware
	r.Use(middleware.AuthMiddleware(user))

	// ws routes
	r.GET("/ws/connect/:uid", wshandler.WsConnector) // * for websocket connection mothod must be GET

	// msg routes
	r.POST("/msgs", api.MakeHTTPHandleFunc(msgHandler.PullAllMessages))
	r.POST("/msgs/group", api.MakeHTTPHandleFunc(msgHandler.PullAllGroupMessages))
	r.PATCH("/msgs", api.MakeHTTPHandleFunc(msgHandler.UpdateIsReadMessage))
	r.DELETE("/msgs", api.MakeHTTPHandleFunc(msgHandler.DeleteMessage))
	r.DELETE("/msgs/group", api.MakeHTTPHandleFunc(msgHandler.DeleteGroupMessage))

	// file routes
	r.POST("/file", api.MakeHTTPHandleFunc(uploadHandler.UploaFile))
	r.DELETE("/file", api.MakeHTTPHandleFunc(uploadHandler.DeleteFile))

	// contact routes
	r.POST("/contacts", api.MakeHTTPHandleFunc(contact.AddContact))
	r.GET("/contacts/:id", api.MakeHTTPHandleFunc(contact.GetContacts))

	// group routes
	r.POST("/groups", api.MakeHTTPHandleFunc(groupHand.CreateGroup))
	r.POST("/groups/addmember", api.MakeHTTPHandleFunc(groupHand.AddMemberToGroup))
	r.GET("/groups/user/:uid", api.MakeHTTPHandleFunc(groupHand.GetAllGroupByUserID))
	r.GET("/groups/:gid", api.MakeHTTPHandleFunc(groupHand.GetGroupDetailsByID))
	r.DELETE("/groups/:gid/:uid", api.MakeHTTPHandleFunc(groupHand.RemoveMemberFromGroup))
	r.DELETE("/groups/:gid", api.MakeHTTPHandleFunc(groupHand.DeleteGroupByID))
	r.PUT("/groups", api.MakeHTTPHandleFunc(groupHand.UpdateGroupDetails))

}

func RunServer(listenAddr string) error {
	return r.RunTLS(listenAddr, "../server.crt", "../server.key")
}
