package main

import (
	"fmt"
	"log"
	"os"

	routes "github.com/floraorg/sakura/router"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := routes.SetupRouter()
	fmt.Printf("Server starting on port :%s\n", port)
	r.Run(":" + port)
}
