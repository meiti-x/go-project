package app

import (
	"fmt"

	"agentic/commerce/config"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

const MyAppName string = "goSocial"

func Banner() string {
	return "\n" + fmt.Sprintf("Service Name: %s", MyAppName)
}

type Application struct {
	Config   *config.Config
	Database *gorm.DB
}

func NewApplication(cfg *config.Config, database *gorm.DB) *Application {
	return &Application{
		Config:   cfg,
		Database: database,
	}
}

var ApplicationModule = fx.Module(
	"application",
	fx.Provide(NewApplication),
)

func (app *Application) Shutdown() {
	app.ShutdownDatabase()
}
func (app *Application) ShutdownDatabase() {
	db, err := app.Database.DB()
	if err == nil {
		_ = db.Close()
	}
}
