package sql

import "database/sql"

type DB struct {
	*sql.DB
}

type Config struct {
	Debug bool
}
