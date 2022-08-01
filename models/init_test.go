package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"testing"
)

func TestConn(t *testing.T) {
	dsn := "host=localhost user=postgres password=123123 dbname=testdb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("gorm Init Error : ", err)
	}
	db.AutoMigrate(&Document{}, &Term{})
	fmt.Println(db.Name())
}
