package app

import (
	"database/sql"
	"embed"

	"github.com/ivangago06/config"
)

type App struct {
	Config         config.Configuration
	Db             *sql.DB
	Res            *embed.FS
	ScheduledTasks Scheduled
}
