package main

import (
	"github.com/google/uuid"
)

type IEntity interface {
    GetId() uuid.UUID
}

type Entity struct {
    Id uuid.UUID `json:"id"`
}

type Profile struct {
    Entity
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (p Entity) GetId() uuid.UUID{
    return p.Id
}

func NewProfile() *Profile {
    return &Profile{}
}

