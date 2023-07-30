package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type ClientInfo struct {
	Name string
	Pwd  string
}

func Register(name string, pwd string) {
	dsn := "root:admin@tcp(127.0.0.1:3306)/p2p?charset=utf8mb4&parseTime=True&loc=Local"
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	client := ClientInfo{
		Name: name,
		Pwd:  pwd,
	}
	db.AutoMigrate(&ClientInfo{})
	db.Create(&client)
}
func Login(name string, pwd string) bool {
	dsn := "root:admin@tcp(127.0.0.1:3306)/p2p?charset=utf8mb4&parseTime=True&loc=Local"
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	client := ClientInfo{
		Name: name,
		Pwd:  pwd,
	}
	tx := db.Where("name = ? AND pwd =?", name, pwd).First(&client)

	if tx.Error != nil {
		log.Println(tx.Error)
		return false
	}
	return true
}
