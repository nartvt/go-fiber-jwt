package database

import (
	"go-fiber-jwt/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	conn, err := gorm.Open(mysql.Open("root:root@/jwt_auth"), &gorm.Config{})
	if err != nil {
		panic("Could not  connect to the database")
	}
	DB = conn
	err1 := conn.AutoMigrate(&models.User{})
	if err1 != nil {
		panic("auto migrate data has an error")
	}
}
