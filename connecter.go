package connecter

type DriverName string

func (d DriverName) String() string {
	return string(d)
}

const (
	DriverNameOfMySQL  DriverName = "mysql"
	DriverNameOfSQLite DriverName = "sqlite"
)
