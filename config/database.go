package config

import (
	"MyGram/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
    host     = "127.0.0.1"
    username = "root"
    password = ""
    port     = 3306
    dbname   = "mygram"
    Database *gorm.DB
    err      error
) 

func Connect() {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
    Database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    
    if err != nil {
        panic(err)
    }

    Database.AutoMigrate(models.Comment{}, models.Photo{}, models.SocialMedia{}, models.User{})
}