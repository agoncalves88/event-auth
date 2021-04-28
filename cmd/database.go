package main

import (
	"fmt"
	"log"

	auth "github.com/agoncalves88/event-auth"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func GetDatabase() *gorm.DB {
	databaseurl := "sqlserver://sa:zaq1ZAQ!@localhost:1433?database=EventDB"

	connection, err := gorm.Open(sqlserver.Open(databaseurl), &gorm.Config{})
	if err != nil {
		log.Fatalln("Invalid database url")
	}
	sqldb, err := connection.DB()
	if err != nil {
		log.Fatal("Error com connect database")
	}
	err = sqldb.Ping()
	if err != nil {
		log.Fatal("Database connected")
	}
	fmt.Println("Database connection successful.")
	return connection
}

func CloseDatabase(connection *gorm.DB) {
	sqldb, _ := connection.DB()
	sqldb.Close()
}

func InitialMigration() {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.AutoMigrate(auth.User{})
}
