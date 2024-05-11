package main

import (
	"github.com/bradleyevansx/inmemory/app"
	"github.com/bradleyevansx/inmemory/bus"
	"github.com/bradleyevansx/inmemory/stor"
)


func main() {
	server := app.NewServer(":3000")
	db, err := stor.ConnectDB()
	if err != nil {
		return
	}
	repo := stor.NewRepository(db)

	repository := bus.NewProfileService(repo)
	controller := app.NewAPIController(server, repository)
	controller.RegisterRoutes()
	server.Run()
}