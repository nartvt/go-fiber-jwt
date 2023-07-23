package infra

import (
	"fmt"
	"go-fiber-jwt/app/models"
	"go-fiber-jwt/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	if DB == nil {
		connect()
	}
}

func connect() {
	conf := config.GetConfig()
	connect := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Mysql.UserName,
		conf.Mysql.Password,
		conf.Mysql.Host,
		conf.Mysql.Port,
		conf.Mysql.Database)
	conn, err := gorm.Open(mysql.Open(connect), &gorm.Config{})
	if err != nil {
		panic("Could not  connect to the database")
	}
	DB = conn
	err1 := conn.AutoMigrate(&models.User{})
	if err1 != nil {
		panic("auto migrate data has an error")
	}
}

func CloseConnect() {
	if DB != nil {
		DB = nil
	}
}
