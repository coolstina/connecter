package mysql

import (
	"time"

	"github.com/coolstina/connecter"

	"gorm.io/gorm/logger"
)

// Config defines config for database.
type Config struct {
	Host                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	LogLevel              int
	Logger                logger.Interface
	DriverName            connecter.DriverName
}
