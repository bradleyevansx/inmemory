package app

import (
	"github.com/bradleyevansx/inmemory/bus"
	"github.com/bradleyevansx/inmemory/stor"
)

type AuthController struct {
	server *Server
	profileRepo bus.IService[stor.Profile]
}

func NewAuthController(server *Server, profileRepo bus.IService[stor.Profile] ) *AuthController{
	return &AuthController{}
}