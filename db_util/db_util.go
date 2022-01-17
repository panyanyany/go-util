package db_util

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDb(name, username, pass string, models []interface{}) *gorm.DB {

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, pass, name)), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(models...)
	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(10)
	return db
}
