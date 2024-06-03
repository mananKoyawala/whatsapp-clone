package main

import (
	"github.com/mananKoyawala/whatsapp-clone/router"
	"log"
)

func main() {
	router.SetupRouters()
	log.Fatal(router.RunServer("localhost:8080"))
}