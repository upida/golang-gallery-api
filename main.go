package main

import (
	"gallery/models"

	"gallery/router"
)

func main() {

	models.ConnectDataBase()

	r := router.SetupRouter()
	router.SetupRoutes(r)

	r.Run(":8888")
}
