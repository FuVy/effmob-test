package main

import (
	"github.com/fuvy/effmob-test/internal/database/migrate"
	"github.com/fuvy/effmob-test/internal/http/handler"
	"github.com/fuvy/effmob-test/internal/service"
	"github.com/fuvy/effmob-test/pkg/env"
	someapi "github.com/fuvy/effmob-test/pkg/someAPI"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Effmob API
// @version 1.0
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.WithError(err).Fatal("no .env file")
	}
	dsn := env.MakePostgresDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.WithError(err).Fatal("failed to connect to database")
	}
	address := env.StringWithDefault("ADDRESS", ":8080")
	someapiUrl := env.GetOrPanicOnEmpty("SOMEAPI_ADDRESS")
	someapi.SetUrl(someapiUrl)
	dbConn, err := db.DB()
	if err != nil {
		log.WithError(err).Fatal("failed to get DBConn")
	}
	if err := migrate.MigrateDB(db); err != nil {
		log.WithError(err).Fatal("failed to migrate db")
	}
	defer dbConn.Close()

	srv := service.New(db)
	handler := handler.New(srv)

	app := fiber.New()
	if err := handler.BindRoutes(app); err != nil {
		log.WithError(err).Fatal("failed to bind routes")
	}
	app.Listen(address)
}
