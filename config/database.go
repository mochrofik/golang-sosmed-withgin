package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	connectString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&loc=Asia%vJakarta",
		ENV.DB_USER, ENV.DB_PASSWORD, ENV.DB_HOST, ENV.DB_PORT, ENV.DB_DATABASE, "%2F")

	db, err := gorm.Open(mysql.Open(connectString), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	DB = db

	log.Println("database connected")
}
