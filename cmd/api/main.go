package main

import (
	"log"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/di"
)

func main() {

	cfg, err := config.LoadConfig(

	if err != nil {
		log.Fatal("Error to load the config")
	}

	server, err := di.InitializeApi(cfg)
	if err != nil {
		log.Fatal("Faild to start the server")
	}

	server.Start()
}
