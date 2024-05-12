package main

import (
	"github.com/bradleyevansx/inmemory/app"
	"github.com/bradleyevansx/inmemory/bus"
	"github.com/bradleyevansx/inmemory/stor"
)
func  newProfileMapper(scan func(dest ...any) error) (*stor.Profile, error) {
	var p stor.Profile
    err := scan(&p.Id, &p.Email, &p.Password, &p.TestInt)
    if err != nil {
        return nil, err
    }
    
    return &p, err
}

func main() {
	server := app.NewServer(":3000")
	db, err := stor.ConnectDB()
	if err != nil {
		return
	}
	repo := stor.NewRepository(db)

	profileService := bus.NewBaseService(repo, "profile", newProfileMapper)
	controller := app.NewAPIController(server, profileService)
	controller.RegisterRoutes()
	server.Run()
}