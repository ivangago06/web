package models

import (
	"time"

	"github.com/ivangago06/app"
	"github.com/ivangago06/database"
)

func RunAllMigratios(app *app.App) error {

	user := User{
		Id:        1,
		UserName:  "migrate",
		Password:  "migrate",
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	}

	err := database.Migrate(app, user)

	if err != nil {
		return err
	}

	session := Session{
		Id:         1,
		UserId:     1,
		AuthToken:  "migrate",
		RememberMe: false,
		CreatedAt:  time.Now(),
	}

	err = database.Migrate(app, session)

	if err != nil {
		return err
	}

	return nil

}
