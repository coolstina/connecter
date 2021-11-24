# database

Quick fast connection to database use gorm.

## Installation

```shell script
$ go get -u github.com/coolstina/connecter
```

## Example 

```go
package main

import (
	"fmt"

	"github.com/coolstina/connecter"
)

func main() {
	db, err := database.New(&database.Options{
		Host:                  "127.0.0.1",
		Username:              "root",
		Password:              "root",
		Database:              "hello",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: 360,
		LogLevel:              1,
		Logger:                nil,
		DriverName:            "mysql",
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("db instance pointer: %+v\n", db)
}
```