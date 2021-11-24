package main

import (
	"fmt"

	"github.com/coolstina/connecter"

	"github.com/coolstina/connecter/mysql"
)

func main() {
	db, err := mysql.NewMySQLConnection(&mysql.Config{
		Host:                  "127.0.0.1",
		Username:              "root",
		Password:              "root",
		Database:              "hello",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: 360,
		LogLevel:              1,
		Logger:                nil,
		DriverName:            connecter.DriverNameOfMySQL,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("db instance pointer: %+v\n", db)
}
