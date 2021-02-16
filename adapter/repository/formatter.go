package repository

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/inhies/go-bytesize"
)

var (
	ErrSingleRecordOverflowsLimits = errors.New("limits are too strict: single record overflows them")
)

type CsvFormatter struct {
	inStream         <-chan []string
	outStream        chan<- io.ReadCloser
	recordsProcessed uint
	writer           *csv.Writer
	buffer           *LimitBuffer
}

func NewCsvFormatter(
	inStream <-chan []string,
	outStream chan<- io.ReadCloser,
	sizeLimit bytesize.ByteSize,
	lineLimit uint,
) *CsvFormatter {
	buffer := NewLimitBuffer(sizeLimit, lineLimit)
	return &CsvFormatter{
		inStream:  inStream,
		outStream: outStream,
		buffer:    buffer,
		writer:    csv.NewWriter(buffer),
	}
}

func (f *CsvFormatter) FormatFiles(ctx context.Context) error {
	for {
		select {
		case record, isOpen := <-f.inStream:
			if !isOpen {
				return f.sendBufferToStream()
			}
			err := f.writeCSVToBuffer(record)
			if err == ErrLinesOverflow || err == ErrSizeOverflow {
				if f.buffer.LinesWritten() == 0 {
					return ErrSingleRecordOverflowsLimits
				} else if err := f.sendBufferToStream(); err != nil {
					return err
				} else if err := f.writeCSVToBuffer(record); err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (f *CsvFormatter) sendBufferToStream() error {
	file, err := f.flushBufferToFile()
	if err != nil {
		return err
	}
	f.outStream <- file
	return nil
}

func (f *CsvFormatter) writeCSVToBuffer(record []string) error {
	if err := f.writer.Write(record); err != nil {
		return err
	}
	f.writer.Flush()
	return f.writer.Error()
}

func (f *CsvFormatter) flushBufferToFile() (io.ReadCloser, error) {
	file, err := os.Create(path.Join("/tmp", uuid.NewString()+".csv"))
	if err != nil {
		return nil, err
	}

	if _, err := file.ReadFrom(f.buffer); err != nil {
		return nil, err
	}
	f.buffer.Reset()
	f.writer = csv.NewWriter(f.buffer)
	if _, err := file.Seek(0, 0); err != nil {
		return nil, err
	}

	return file, nil
}
