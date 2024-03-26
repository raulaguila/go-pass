package database

import (
	"context"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/raulaguila/go-pass/internal/pkg/domain"
	"github.com/raulaguila/go-pass/pkg/helpers"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func autoMigrate(postgresdb *gorm.DB) {
	helpers.PanicIfErr(postgresdb.AutoMigrate(&domain.Permissions{}))
	helpers.PanicIfErr(postgresdb.AutoMigrate(&domain.Profile{}))
	helpers.PanicIfErr(postgresdb.AutoMigrate(&domain.User{}))

	helpers.PanicIfErr(postgresdb.AutoMigrate(&domain.Operator{}))
	helpers.PanicIfErr(postgresdb.AutoMigrate(&domain.Phone{}))
	helpers.PanicIfErr(postgresdb.AutoMigrate(&domain.Site{}))
	helpers.PanicIfErr(postgresdb.AutoMigrate(&domain.Account{}))
	helpers.PanicIfErr(postgresdb.AutoMigrate(&domain.AccountMailHistory{}))
}

func createDefaults(postgresdb *gorm.DB) {
	profiles := []domain.Profile{
		{
			Name: "ROOT",
			Permissions: domain.Permissions{
				UserModule:    true,
				ProfileModule: true,
			},
		}, {
			Name: "USER",
			Permissions: domain.Permissions{
				UserModule:    false,
				ProfileModule: false,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var profileID uint = 0
	for i, profile := range profiles {
		helpers.PanicIfErr(postgresdb.WithContext(ctx).FirstOrCreate(&profile, "name = ?", profile.Name).Error)
		if i == 0 {
			profileID = profile.Id
		}
	}

	user := &domain.User{
		Name:      os.Getenv("ADM_NAME"),
		Email:     os.Getenv("ADM_MAIL"),
		Status:    true,
		ProfileID: profileID,
		New:       false,
		Token:     new(string),
		Password:  new(string),
	}

	token := uuid.New().String()
	*user.Token = token

	hash, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADM_PASS")), bcrypt.DefaultCost)
	helpers.PanicIfErr(err)
	user.Password = new(string)
	*user.Password = string(hash)

	helpers.PanicIfErr(postgresdb.WithContext(ctx).FirstOrCreate(user, "mail = ?", user.Email).Error)
}
