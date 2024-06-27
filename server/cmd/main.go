package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mananKoyawala/whatsapp-clone/config"
	db "github.com/mananKoyawala/whatsapp-clone/database"
	"github.com/mananKoyawala/whatsapp-clone/router"
)

func main() {

	// load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	// get AWS s3 variables
	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_KEY")
	bucketName := os.Getenv("AWS_BUCKET_NAME")

	// get AES encryption secrect key and iv
	encSecretKey := os.Getenv("ENC_SECRET_KEY")
	encIv := os.Getenv("ENC_INITIALIZATION_VECTOR")

	// database connection
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err.Error())
	}

	// get Configuration
	hub := config.Configuration(db.GetDB(), region, bucketName, accessKey, secretKey, encSecretKey, encIv)

	//run the hub
	go hub.Run()

	log.Fatal(router.RunServer("localhost:8080"))
}
