package repository

import (
	"context"
	"database/sql"
)

type (
	Database interface {
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	}

	SqlDataFetcher struct {
		OutStream   chan<- []string
		CountQuery  string
		SelectQuery string
		Db          Database
	}
)

func (s *SqlDataFetcher) CountRecords(ctx context.Context) (uint, error) {
	row := s.Db.QueryRowContext(ctx, s.CountQuery)
	if row.Err() != nil {
		return 0, row.Err()
	}
	var res uint
	if err := row.Scan(&res); err != nil {
		return 0, err
	}
	return res, nil
}

func (s *SqlDataFetcher) StreamData(ctx context.Context) error {
	rows, err := s.Db.QueryContext(ctx, s.SelectQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	s.OutStream <- cols

	values := make([]interface{}, len(cols))
	for i := range cols {
		values[i] = new(sql.RawBytes)
	}

	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return err
		}
		s.OutStream <- rawBytesToString(values)
	}
	return rows.Err()
}

func rawBytesToString(values []interface{}) []string {
	res := make([]string, len(values))
	for idx, v := range values {
		res[idx] = string(*v.(*sql.RawBytes))
	}
	return res
}
