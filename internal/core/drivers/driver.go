package drivers

type Motor string

const (
	SQLite     Motor = "sqlite"
	MariaDB    Motor = "mariadb"
	PostgreSQL Motor = "postgresql"
)

type Driver struct {
	Motor        Motor
	GetType      func(tryType string) string
	GetAttribute func(tryAttribute string) string
	GetFunction  func(tryFunction string) string
	GetLiteral   func(tryLiteral string) string
}
