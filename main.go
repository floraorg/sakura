package main

import (
	"log"

	routes "github.com/floraorg/sakura/router"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := routes.SetupRouter()
	r.Run()
}
