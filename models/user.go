package models

import (
	"log"
	"net/http"

	//"os/user"
	"time"

	"github.com/ivangago06/app"
)

type User struct {
	Id        int64
	UserName  string
	Password  string
	CreatedAt time.Time
	UpdateAt  time.Time
}

const userColumnsNoId = "\"Username\", \"Password\", \"CreatedAt\", \"UpdatedAt\""
const userColumns = "\"Id\", " + userColumnsNoId
const userTable = "public.\"User\""

const (
	selectUserById       = "SELECT " + userColumns + " FROM " + userTable + " WHERE \"Id\" = $1"
	selectUserByUsername = "SELECT " + userColumns + " FROM " + userTable + " WHERE \"Username\" = $1"
	insertUser           = "INSERT INTO " + userTable + " (" + userColumnsNoId + ") VALUES ($1, $2, $3, $4) RETURNING \"Id\""
)

func GetCurrentUser(app *app.App, r *http.Request) (User, error) {
	cookie, err := r.Cookie("session")

	if err != nil {
		log.Println("No se encontró la sesión.")
		return User{}, err
	}

	session, err := GetSessionByAuthToken(app, cookie.Value)

	return GetUserById(app, session)
}
