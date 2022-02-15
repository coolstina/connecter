// Copyright 2021 helloshaohua <wu.shaohua@foxmail.com>;
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

package mysql

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewConnection create a new gorm db instance with the given options.
func NewConnection(config *Config, ops ...Option) (*gorm.DB, error) {

	// If not exists then create.
	err := CreateDatabaseIfNotExists(
		config.DriverName.String(),
		NewDataSourceNameForNoSelectDatabase(config.Host, config.Username, config.Password, ops...),
		config.Database,
	)

	if err != nil {
		return nil, err
	}

	dsn := NewDataSourceNameForConfig(config, ops...)
	opts := &gorm.Config{Logger: config.Logger}

	db, err := gorm.Open(mysql.Open(dsn), opts)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(config.MaxOpenConnections)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(config.MaxConnectionLifeTime)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(config.MaxIdleConnections)

	return db, nil
}

// NewDataSourceNameForConfig Get data source name for given options.
func NewDataSourceNameForConfig(config *Config, ops ...Option) string {
	return NewDataSourceName(config.Host, config.Username, config.Password, config.Database, ops...)
}

func NewDataSourceNameForNoSelectDatabase(host, username, password string, ops ...Option) string {
	return NewDataSourceName(host, username, password, "")
}

// NewDataSourceName initialize database of the data source name.
// If database parameter is empty, will not choose database, such as only open database connection.
func NewDataSourceName(host, username, password, database string, ops ...Option) string {
	options := &options{
		charset:   "utf8mb4",
		parseTime: true,
		location:  "Local",
	}

	for _, o := range ops {
		o.apply(options)
	}

	args := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
		username,
		password,
		host,
		database,
		options.charset,
		options.parseTime,
		options.location,
	)

	return args
}

// CreateDatabaseIfNotExists If database not exists, then create it.
func CreateDatabaseIfNotExists(driverName, dataSourceName, databaseName string, ops ...Option) error {
	options := &options{
		charset: "utf8mb4",
	}

	for _, o := range ops {
		o.apply(options)
	}

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}
	defer db.Close()

	s := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET %s", databaseName, options.charset)
	if _, err = db.Exec(s); err != nil {
		return err
	}

	return nil
}

// GormComment get table comment for the description.
func GormComment(description string) string {
	return fmt.Sprintf("comment '%s'", description)
}
