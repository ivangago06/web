package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ivangago06/app"
)

// función que hace la conexión a la BD

func ConnectDB(app *app.App) *sql.DB {

	postgresonfig := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", app.Config.Db.Ip, app.Config.Db.Port, app.Config.Db.User, app.Config.Db.Password, app.Config.Db.Name)

	db, err := sql.Open("postgres", postgresonfig)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	log.Println("Conexónxitosa a neustra BD.")

	return db

}
