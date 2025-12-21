package driver

type Motor string

const (
	SQLite   Motor = "sqlite"
	MariaDB  Motor = "mariadb"
	MySQL    Motor = "mysql"
	Postgres Motor = "postgres"
)

type Driver struct {
	Motor      Motor
	Types      Types
	Variations Variations
}
