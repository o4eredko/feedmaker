package repository_test

import (
	"context"
	"database/sql/driver"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"go-feedmaker/adapter/repository"
	helper "go-feedmaker/infrastructure/testing"
)

type sqlFetcherFields struct {
	CountQuery  string
	SelectQuery string
}

func TestSqlDataFetcher_StreamData(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	testCases := []struct {
		name         string
		args         *args
		fields       *sqlFetcherFields
		inCsvRecords [][]string
		setupMocks   func(*args, *sqlFetcherFields, sqlmock.Sqlmock, [][]string)
		wantErr      error
	}{
		{
			name: "succeed",
			args: &args{ctx: context.Background()},
			fields: &sqlFetcherFields{
				SelectQuery: "SELECT * FROM Marketing.dbo.records;",
				CountQuery:  "SELECT Count(*) FROM Marketing.dbo.records;",
			},
			inCsvRecords: helper.ReadCsvFromFile(t, "testdata/records.csv"),
			setupMocks: func(args *args, f *sqlFetcherFields, sql sqlmock.Sqlmock, csvRecords [][]string) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(100)
				sql.ExpectQuery(regexp.QuoteMeta(f.CountQuery)).WillReturnRows(rows).RowsWillBeClosed()

				rows = sqlmock.NewRows([]string{"col1", "col2", "col3", "col4"})
				for _, record := range csvRecords {
					rows.AddRow(csvRecordToSqlValues(record)...)
				}
				sql.ExpectQuery(regexp.QuoteMeta(f.SelectQuery)).RowsWillBeClosed().WillReturnRows(rows)
			},
		},
		{
			name: "Count Query error",
			args: &args{ctx: context.Background()},
			fields: &sqlFetcherFields{
				SelectQuery: regexp.QuoteMeta("SELECT * FROM Marketing.dbo.records;"),
				CountQuery:  regexp.QuoteMeta("SELECT Count(*) FROM Marketing.dbo.records;"),
			},
			inCsvRecords: helper.ReadCsvFromFile(t, "testdata/records.csv"),
			setupMocks: func(args *args, f *sqlFetcherFields, sql sqlmock.Sqlmock, csvRecords [][]string) {
				sql.ExpectQuery(regexp.QuoteMeta(f.CountQuery)).WillReturnError(defaultErr)
			},
			wantErr: defaultErr,
		},
		{
			name: "Select Query error",
			args: &args{ctx: context.Background()},
			fields: &sqlFetcherFields{
				SelectQuery: regexp.QuoteMeta("SELECT * FROM Marketing.dbo.records;"),
				CountQuery:  regexp.QuoteMeta("SELECT Count(*) FROM Marketing.dbo.records;"),
			},
			inCsvRecords: helper.ReadCsvFromFile(t, "testdata/records.csv"),
			setupMocks: func(args *args, f *sqlFetcherFields, sql sqlmock.Sqlmock, csvRecords [][]string) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(100)
				sql.ExpectQuery(regexp.QuoteMeta(f.CountQuery)).WillReturnRows(rows).RowsWillBeClosed()
				sql.ExpectQuery(regexp.QuoteMeta(f.SelectQuery)).WillReturnError(defaultErr)
			},
			wantErr: defaultErr,
		},
		{
			name: "context error",
			args: &args{ctx: helper.TimeoutCtx(t, context.Background(), time.Nanosecond)},
			fields: &sqlFetcherFields{
				SelectQuery: "SELECT * FROM Marketing.dbo.records;",
				CountQuery:  "SELECT Count(*) FROM Marketing.dbo.records;",
			},
			inCsvRecords: helper.ReadCsvFromFile(t, "testdata/records.csv"),
			setupMocks: func(args *args, f *sqlFetcherFields, sql sqlmock.Sqlmock, csvRecords [][]string) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(100)
				sql.ExpectQuery(regexp.QuoteMeta(f.CountQuery)).WillReturnRows(rows).WillDelayFor(time.Second)
			},
			wantErr: sqlmock.ErrCancelled,
		},
		{
			name: "row error",
			args: &args{ctx: context.Background()},
			fields: &sqlFetcherFields{
				SelectQuery: "SELECT * FROM Marketing.dbo.records;",
				CountQuery:  "SELECT Count(*) FROM Marketing.dbo.records;",
			},
			inCsvRecords: helper.ReadCsvFromFile(t, "testdata/records.csv"),
			setupMocks: func(args *args, f *sqlFetcherFields, sql sqlmock.Sqlmock, csvRecords [][]string) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(100)
				sql.ExpectQuery(regexp.QuoteMeta(f.CountQuery)).WillReturnRows(rows).RowsWillBeClosed()

				rows = sqlmock.NewRows([]string{"col1", "col2", "col3", "col4"})
				for _, record := range csvRecords {
					rows.AddRow(csvRecordToSqlValues(record)...)
				}
				rows.RowError(5, defaultErr)
				sql.ExpectQuery(regexp.QuoteMeta(f.SelectQuery)).RowsWillBeClosed().WillReturnRows(rows)
			},
			wantErr: defaultErr,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			conn, err := db.Conn(context.Background())
			assert.NoError(t, err)
			defer assert.NoError(t, db.Close())

			recordStream := make(chan []string, 10)
			tc.setupMocks(tc.args, tc.fields, mock, tc.inCsvRecords)

			sqlFetcher := repository.SqlDataFetcher{
				Db:          conn,
				OutStream:   recordStream,
				SelectQuery: tc.fields.SelectQuery,
				CountQuery:  tc.fields.CountQuery,
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				gotErr := sqlFetcher.StreamData(tc.args.ctx)
				assert.Equal(t, tc.wantErr, gotErr)
				close(recordStream)
				wg.Done()
			}()

			if tc.wantErr == nil {
				var idx int
				assert.Equal(t, []string{"col1", "col2", "col3", "col4"}, <-recordStream)
				for record := range recordStream {
					assert.Equal(t, tc.inCsvRecords[idx], record)
					idx++
				}
			}
			wg.Wait()
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func csvRecordToSqlValues(record []string) []driver.Value {
	values := make([]driver.Value, 0, len(record))
	for _, value := range record {
		values = append(values, value)
	}
	return values
}
