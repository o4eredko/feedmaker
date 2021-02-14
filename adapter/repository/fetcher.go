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
		OutStream        chan<- []string
		CountQuery       string
		SelectQuery      string
		Db               Database
		recordsCount     uint
		recordsProceeded uint
		progress         uint
		onDataFetched    func()
		onProgress       func(progress uint)
	}
)

func (s *SqlDataFetcher) countRecords(ctx context.Context) error {
	row := s.Db.QueryRowContext(ctx, s.CountQuery)
	if row.Err() != nil {
		return row.Err()
	}
	if err := row.Scan(&s.recordsCount); err != nil {
		return err
	}
	return nil
}

func (s *SqlDataFetcher) StreamData(ctx context.Context) error {
	if err := s.countRecords(ctx); err != nil {
		return err
	}
	rows, err := s.Db.QueryContext(ctx, s.SelectQuery)
	if err != nil {
		return err
	}
	defer rows.Close()
	if s.onDataFetched != nil {
		s.onDataFetched()
	}

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
		s.recordsProceeded++
		s.updateProgress()
	}
	return rows.Err()
}

func (s *SqlDataFetcher) updateProgress() {
	if s.recordsProceeded/s.recordsCount > s.progress {
		s.progress = s.recordsProceeded / s.recordsCount
		if s.onProgress != nil {
			s.onProgress(s.progress)
		}
	}
}

func rawBytesToString(values []interface{}) []string {
	res := make([]string, len(values))
	for idx, v := range values {
		res[idx] = string(*v.(*sql.RawBytes))
	}
	return res
}

func (s *SqlDataFetcher) OnDataFetched(callback func()) {
	s.onDataFetched = callback
}

func (s *SqlDataFetcher) OnProgress(callback func(progress uint)) {
	s.onProgress = callback
}
