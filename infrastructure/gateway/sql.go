package gateway

import (
	"database/sql"
	"errors"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	ErrSqlDisconnected = errors.New("gateway is not connected to Sql")
)

type SqlGateway struct {
	DriverName string
	DSN        string
	db         *sql.DB
}

func (s *SqlGateway) Connect() error {
	db, err := sql.Open(s.DriverName, s.DSN)
	if err != nil {
		return err
	}
	s.db = db
	return db.Ping()
}

func (s *SqlGateway) DB() *sql.DB {
	return s.db
}

func (s *SqlGateway) Disconnect() error {
	if s.db == nil {
		return ErrSqlDisconnected
	}
	return s.db.Close()
}
