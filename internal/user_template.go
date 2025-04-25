package internal

import (
	"time"

	"github.com/junsazanami430u/test-bob/pkg/gen/models/factory"
	"github.com/oklog/ulid/v2"
)

func BobFactryExample(name string, email string, password string) *factory.UserTemplate {
	f := factory.New()

	f.AddBaseUserMod(
		factory.UserMods.IDFunc(
			func() []byte {
				return ulid.Make().Bytes()
			},
		),
		factory.UserMods.Name(name),
		factory.UserMods.Email(email),
		factory.UserMods.Password(password),
		factory.UserMods.CreatedAtFunc(
			func() time.Time {
				return time.Now()
			},
		),
	)

	return f.NewUser()
}
