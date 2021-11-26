// Copyright 2021 helloshaohua &lt;wu.shaohua@foxmail.com&gt;
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	"github.com/coolstina/connecter"

	"github.com/coolstina/connecter/mysql"
)

func main() {
	db, err := mysql.NewConnection(&mysql.Config{
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
