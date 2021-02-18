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
	csvWriter        *csv.Writer
	limitWriter      *LimitWriter
	currentFile      *os.File
	sizeLimit        bytesize.ByteSize
	lineLimit        uint
}

func NewCsvFormatter(
	inStream <-chan []string,
	outStream chan<- io.ReadCloser,
	sizeLimit bytesize.ByteSize,
	lineLimit uint,
) *CsvFormatter {
	return &CsvFormatter{
		inStream:  inStream,
		outStream: outStream,
		sizeLimit: sizeLimit,
		lineLimit: lineLimit,
	}
}

func (f *CsvFormatter) FormatFiles(ctx context.Context) error {
	if err := f.createCsvWriter(); err != nil {
		return err
	}
	for {
		select {
		case record, isOpen := <-f.inStream:
			if !isOpen {
				return f.sendCsvFileToStream()
			}
			err := f.writeRecordToCsv(record)
			if err == ErrLinesOverflow || err == ErrSizeOverflow {
				if f.limitWriter.LinesWritten() == 0 {
					return ErrSingleRecordOverflowsLimits
				}
				if err := f.sendCsvFileToStream(); err != nil {
					return err
				}
				if err := f.createCsvWriter(); err != nil {
					return err
				}
				if err := f.writeRecordToCsv(record); err != nil {
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

func (f *CsvFormatter) createCsvWriter() error {
	file, err := createTmpFile()
	if err != nil {
		return err
	}
	f.currentFile = file
	f.limitWriter = NewLimitWriter(file, f.sizeLimit, f.lineLimit)
	f.csvWriter = csv.NewWriter(f.limitWriter)
	return nil
}

func (f *CsvFormatter) sendCsvFileToStream() error {
	if err := f.limitWriter.Flush(); err != nil {
		return err
	} else if _, err := f.currentFile.Seek(0, 0); err != nil {
		return err
	}
	f.outStream <- f.currentFile
	return nil
}

func (f *CsvFormatter) writeRecordToCsv(record []string) error {
	if err := f.csvWriter.Write(record); err != nil {
		return err
	}
	f.csvWriter.Flush()
	return f.csvWriter.Error()
}

func createTmpFile() (*os.File, error) {
	return os.Create(path.Join("/tmp", uuid.NewString()+".csv"))
}
