package main

import (
	"database/sql"
)

type ProfileService struct {
	*BaseRepository[Profile]
}

func  newProfileMapper(scan func(dest ...any) error) (*Profile, error) {
	var p Profile
    err := scan(&p.Id, &p.Email, &p.Password, &p.TestInt)
    if err != nil {
        return nil, err
    }
    
    return &p, err
}

func NewProfileService(db *sql.DB) *ProfileService {
	return &ProfileService{BaseRepository: NewRepository(db, "profile", newProfileMapper)}
}