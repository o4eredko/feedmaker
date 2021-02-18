package repository_test

import (
	"bufio"
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/inhies/go-bytesize"
	"github.com/stretchr/testify/assert"

	"go-feedmaker/adapter/repository"
	helper "go-feedmaker/infrastructure/testing"
)

type (
	fields struct {
		sizeLimit bytesize.ByteSize
		lineLimit uint
		inStream  chan []string
		outStream chan io.ReadCloser
	}
	record struct {
		in  []string
		csv string
	}
)

func TestCsvFormatter_FormatFiles(t *testing.T) {
	testCases := []struct {
		name    string
		fields  *fields
		ctx     context.Context
		records []record
		wantErr error
	}{
		{
			name: "line limit checked",
			fields: &fields{
				sizeLimit: bytesize.MB,
				lineLimit: 3,
				inStream:  make(chan []string),
				outStream: make(chan io.ReadCloser),
			},
			ctx: context.Background(),
			records: []record{
				{in: []string{"a1", "b1", "c1"}, csv: "a1,b1,c1"},
				{in: []string{"a2", "b2", "c2"}, csv: "a2,b2,c2"},
				{in: []string{"a3", "b3", "c3"}, csv: "a3,b3,c3"},
				{in: []string{"a4", "b4", "c4"}, csv: "a4,b4,c4"},
				{in: []string{"a5", "b5", "c5"}, csv: "a5,b5,c5"},
			},
		},
		{
			name: "size limit checked",
			fields: &fields{
				sizeLimit: 10 * bytesize.B,
				lineLimit: 100,
				inStream:  make(chan []string),
				outStream: make(chan io.ReadCloser),
			},
			ctx: context.Background(),
			records: []record{
				{in: []string{"a1", "b1", "c1"}, csv: "a1,b1,c1"},
				{in: []string{"a2", "b2", "c2"}, csv: "a2,b2,c2"},
				{in: []string{"a3", "b3", "c3"}, csv: "a3,b3,c3"},
				{in: []string{"a4", "b4", "c4"}, csv: "a4,b4,c4"},
				{in: []string{"a5", "b5", "c5"}, csv: "a5,b5,c5"},
			},
		},
		{
			name: "size limit lower than single record",
			fields: &fields{
				sizeLimit: 1,
				lineLimit: 100,
				inStream:  make(chan []string),
				outStream: make(chan io.ReadCloser),
			},
			ctx: context.Background(),
			records: []record{
				{in: []string{"a1", "b1", "c1"}, csv: "a1,b1,c1"},
				{in: []string{"a2", "b2", "c2"}, csv: "a2,b2,c2"},
				{in: []string{"a3", "b3", "c3"}, csv: "a3,b3,c3"},
			},
			wantErr: repository.ErrSingleRecordOverflowsLimits,
		},
		{
			name: "line limit is 0",
			fields: &fields{
				sizeLimit: bytesize.MB,
				lineLimit: 0,
				inStream:  make(chan []string),
				outStream: make(chan io.ReadCloser),
			},
			ctx: context.Background(),
			records: []record{
				{in: []string{"a", "b1", "c1"}, csv: "a1,b1,c1"},
				{in: []string{"a2", "b2", "c2"}, csv: "a2,b2,c2"},
				{in: []string{"a3", "b3", "c3"}, csv: "a3,b3,c3"},
			},
			wantErr: repository.ErrSingleRecordOverflowsLimits,
		},
		{
			name: "context error",
			fields: &fields{
				sizeLimit: bytesize.MB,
				lineLimit: 100,
				inStream:  make(chan []string),
				outStream: make(chan io.ReadCloser),
			},
			ctx: helper.TimeoutCtx(t, context.Background(), 0),
			records: []record{
				{in: []string{"a1", "b1", "c1"}, csv: "a1,b1,c1"},
				{in: []string{"a2", "b2", "c2"}, csv: "a2,b2,c2"},
				{in: []string{"a3", "b3", "c3"}, csv: "a3,b3,c3"},
			},
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			formatter := repository.NewCsvFormatter(
				tc.fields.inStream,
				tc.fields.outStream,
				tc.fields.sizeLimit,
				tc.fields.lineLimit,
			)

			var wg sync.WaitGroup
			wg.Add(2)
			produceRecordsCtx, cancelProducing := context.WithCancel(context.Background())
			defer cancelProducing()
			go func() {
				defer wg.Done()
				defer close(tc.fields.outStream)
				gotErr := formatter.FormatFiles(tc.ctx)
				if gotErr != nil {
					cancelProducing()
				}
				assert.Equal(t, tc.wantErr, gotErr)
			}()
			go func() {
				defer wg.Done()
				defer close(tc.fields.inStream)
				produceRecords(produceRecordsCtx, tc.fields.inStream, tc.records)
			}()

			var linesChecked int
			for file := range tc.fields.outStream {
				var fileSize, fileLineCount int
				scanner := bufio.NewScanner(file)
				scanner.Split(bufio.ScanLines)
				for scanner.Scan() {
					assert.Equal(t, tc.records[linesChecked].csv, scanner.Text())
					fileSize += len(scanner.Bytes())
					fileLineCount++
					linesChecked++
				}
				assert.NoError(t, scanner.Err())
				assert.LessOrEqual(t, uint(fileLineCount), tc.fields.lineLimit)
				assert.LessOrEqual(t, fileSize, int(tc.fields.sizeLimit))
				assert.NoError(t, file.Close())
			}
			wg.Wait()
		})
	}
}

func produceRecords(ctx context.Context, stream chan<- []string, records []record) {
	if len(records) == 0 {
		return
	}
	var idx int
	for idx < len(records) {
		select {
		case <-ctx.Done():
			return
		case stream <- records[idx].in:
			idx++
		default:
			time.Sleep(time.Millisecond)
		}
	}
}
