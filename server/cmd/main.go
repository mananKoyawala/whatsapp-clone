package main

import (
	"log"

	db "github.com/mananKoyawala/whatsapp-clone/database"
	"github.com/mananKoyawala/whatsapp-clone/internal/user"
	"github.com/mananKoyawala/whatsapp-clone/router"
)

func main() {
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("DB connected")
	userRepository := user.NewUserRepository(db.GetDB())
	userService := user.NewUserService(userRepository)
	userHandler := user.NewUserHandler(userService)
	router.SetupRouters(userHandler)
	log.Fatal(router.RunServer("localhost:8080"))
}
