package config

import (
	"database/sql"

	"github.com/mananKoyawala/whatsapp-clone/internal/contact"
	"github.com/mananKoyawala/whatsapp-clone/internal/group"
	msg "github.com/mananKoyawala/whatsapp-clone/internal/message"
	"github.com/mananKoyawala/whatsapp-clone/internal/user"
	"github.com/mananKoyawala/whatsapp-clone/internal/ws"
	logger "github.com/mananKoyawala/whatsapp-clone/logging"
	"github.com/mananKoyawala/whatsapp-clone/router"
	"github.com/mananKoyawala/whatsapp-clone/service/security"
	"github.com/mananKoyawala/whatsapp-clone/service/upload"
)

func Configuration(db *sql.DB, region, bucketName, accessKey, secretKey, aesSecretKey, aesIv string) *ws.Hub {

	// initialize loggers
	userLogger := logger.InitUserLogger()
	messageLogger := logger.InitMessageLogger()
	groupLogger := logger.InitGroupLogger()
	contactLogger := logger.InitContactLogger()
	wsLogger := logger.InitWSLogger()

	// initalize encryption service
	aesEnc := security.NewEncyption(aesSecretKey, aesIv)

	// user
	userRepository := user.NewUserRepository(db, userLogger)
	userService := user.NewUserService(userRepository, userLogger)
	userHandler := user.NewUserHandler(userService, userLogger)

	// group initialization
	groupRepo := group.NewGroupRepository(db, groupLogger)
	groupSev := group.NewGroupService(groupRepo, groupLogger)
	groupHand := group.NewGroupHandler(groupSev, groupLogger)

	// message
	msgRepo := msg.NewMsgReposritory(db, messageLogger, *aesEnc)
	msgSev := msg.NewMsgService(msgRepo, userRepository, groupRepo, messageLogger)
	msgHand := msg.NewMsgHandler(msgSev, messageLogger)

	// ws initialization
	hub := ws.NewHub()
	wsHandler := ws.NewWsHandler(hub, msgRepo, groupRepo, wsLogger)

	// file upload initialization
	uploadSev := upload.NewAwsService(region, accessKey, secretKey, bucketName)
	uploadSev.InitializeAwsSerive(region, accessKey, secretKey)
	uploadHan := upload.NewAwsHandler(*uploadSev)

	// contact initialization
	conRepo := contact.NewContactRepo(db, contactLogger)
	conSev := contact.NewContactServ(conRepo, userRepository, contactLogger)
	conHand := contact.NewContactHan(conSev, contactLogger)

	router.SetupRouters(userHandler, wsHandler, msgHand, &uploadHan, &conHand, &groupHand)

	return hub
}
