package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	db "github.com/mananKoyawala/whatsapp-clone/database"
	"github.com/mananKoyawala/whatsapp-clone/internal/contact"
	"github.com/mananKoyawala/whatsapp-clone/internal/group"
	msg "github.com/mananKoyawala/whatsapp-clone/internal/message"
	"github.com/mananKoyawala/whatsapp-clone/internal/user"
	"github.com/mananKoyawala/whatsapp-clone/internal/ws"
	logger "github.com/mananKoyawala/whatsapp-clone/logging"
	"github.com/mananKoyawala/whatsapp-clone/router"
	"github.com/mananKoyawala/whatsapp-clone/service/upload"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_KEY")
	bucketName := os.Getenv("AWS_BUCKET_NAME")

	db, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err.Error())
	}

	// log.Println("DB connected")

	// initialize loggers
	userLogger := logger.InitUserLogger()
	messageLogger := logger.InitMessageLogger()

	// user
	userRepository := user.NewUserRepository(db.GetDB(), userLogger)
	userService := user.NewUserService(userRepository, userLogger)
	userHandler := user.NewUserHandler(userService, userLogger)

	// group initialization
	groupRepo := group.NewGroupRepository(db.GetDB())
	groupSev := group.NewGroupService(groupRepo)
	groupHand := group.NewGroupHandler(groupSev)

	// message
	msgRepo := msg.NewMsgReposritory(db.GetDB(), messageLogger)
	msgSev := msg.NewMsgService(msgRepo, userRepository, groupRepo, messageLogger)
	msgHand := msg.NewMsgHandler(msgSev, messageLogger)

	// ws initialization
	hub := ws.NewHub()
	wsHandler := ws.NewWsHandler(hub, msgRepo, groupRepo)

	// file upload initialization
	uploadSev := upload.NewAwsService(region, accessKey, secretKey, bucketName)
	uploadSev.InitializeAwsSerive(region, accessKey, secretKey)
	uploadHan := upload.NewAwsHandler(*uploadSev)

	// contact initialization
	conRepo := contact.NewContactRepo(db.GetDB())
	conSev := contact.NewContactServ(conRepo, userRepository)
	conHand := contact.NewContactHan(conSev)

	//run the hub
	go hub.Run()

	router.SetupRouters(userHandler, wsHandler, msgHand, &uploadHan, &conHand, &groupHand)
	log.Fatal(router.RunServer("localhost:8080"))
}
