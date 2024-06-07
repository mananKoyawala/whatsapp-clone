package main

import (
	"log"

	db "github.com/mananKoyawala/whatsapp-clone/database"
	msg "github.com/mananKoyawala/whatsapp-clone/internal/message"
	"github.com/mananKoyawala/whatsapp-clone/internal/user"
	"github.com/mananKoyawala/whatsapp-clone/internal/ws"
	"github.com/mananKoyawala/whatsapp-clone/router"
)

func main() {
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("DB connected")

	// user
	userRepository := user.NewUserRepository(db.GetDB())
	userService := user.NewUserService(userRepository)
	userHandler := user.NewUserHandler(userService)

	// message
	msgRepo := msg.NewMsgReposritory(db.GetDB())
	msgSev := msg.NewMsgService(msgRepo, userRepository)
	msgHand := msg.NewMsgHandler(msgSev)

	// ws initalization
	hub := ws.NewHub()
	wsHandler := ws.NewWsHandler(hub)

	//run the hub
	go hub.Run()

	router.SetupRouters(userHandler, wsHandler, msgHand)
	log.Fatal(router.RunServer("localhost:8080"))
}
