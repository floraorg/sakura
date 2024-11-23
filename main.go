package main

import (
	"github.com/floraorg/sakura/router"
)

func main() {
	r := routes.SetupRouter()
	r.Run()
}