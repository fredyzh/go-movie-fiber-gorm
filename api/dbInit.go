package api

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (app *Application) connectToDB() (*gorm.DB, error) {
	log.Println(app.DSN)
	db, err := gorm.Open(postgres.Open(app.DSN), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: true,
	})

	if err != nil {
		log.Println("Failed to connect to database: ", err)
		return nil, err
	}

	log.Println("db connected")
	db.Logger = logger.Default.LogMode(logger.Silent)

	return db, nil
}
