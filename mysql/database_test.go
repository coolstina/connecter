package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var def = &Config{
	Host:                  "127.0.0.1",
	Username:              "root",
	Password:              "root",
	Database:              "shaohua4",
	MaxIdleConnections:    100,
	MaxOpenConnections:    100,
	MaxConnectionLifeTime: 360,
	LogLevel:              1,
	Logger:                nil,
	DriverName:            "mysql",
}

func TestNewDataSourceName_NoneSelectDatabase(t *testing.T) {
	expected := `root:root@tcp(127.0.0.1)/?charset=utf8mb4&parseTime=true&loc=Local`
	actual := NewDataSourceName(def.Host, def.Username, def.Password, "")
	assert.Equal(t, expected, actual)
}

func TestNewDataSourceName_SelectDatabase(t *testing.T) {
	expected := `root:root@tcp(127.0.0.1)/mysql?charset=utf8mb4&parseTime=true&loc=Local`
	actual := NewDataSourceName(def.Host, def.Username, def.Password, "mysql")
	assert.Equal(t, expected, actual)
}

func TestCreateDatabaseIfNotExists(t *testing.T) {
	dsn := NewDataSourceName(def.Host, def.Username, def.Password, "")
	assert.NotEmpty(t, dsn)

	err := CreateDatabaseIfNotExists(def.DriverName.String(), dsn, def.Database)
	assert.NoError(t, err)
}

func TestCreateDatabaseIfNotExists_WithCharset(t *testing.T) {
	ops := []Option{
		WithCharset("utf8mb4"),
	}

	dsn := NewDataSourceName(def.Host, def.Username, def.Password, "", ops...)
	assert.NotEmpty(t, dsn)

	err := CreateDatabaseIfNotExists(def.DriverName.String(), dsn, def.Database, ops...)
	assert.NoError(t, err)
}

func TestNewMySQLConnection(t *testing.T) {
	connection, err := NewMySQLConnection(def)
	assert.NoError(t, err)
	assert.NotNil(t, connection)
}
