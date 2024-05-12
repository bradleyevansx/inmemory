package bus

// import (
// 	"github.com/bradleyevansx/inmemory/stor"
// )

// type ProfileService struct {
// 	*BaseService[stor.Profile]
// }

// func  newProfileMapper(scan func(dest ...any) error) (*stor.Profile, error) {
// 	var p stor.Profile
//     err := scan(&p.Id, &p.Email, &p.Password, &p.TestInt)
//     if err != nil {
//         return nil, err
//     }

//     return &p, err
// }

// func NewProfileService(repo *stor.Repository) *ProfileService {
// 	return &ProfileService{BaseService: NewBaseService(repo, "profile", newProfileMapper)}
// }