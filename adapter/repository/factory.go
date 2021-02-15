package repository

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/inhies/go-bytesize"

	"go-feedmaker/interactor"
)

type (
	defaultFactory struct {
		fileSizeLimit  bytesize.ByteSize
		fileLineLimit  uint
		generationType string
		countQuery     string
		selectQuery    string
		sqlGateway     SqlGateway
		ftpGateway     FtpGateway
	}
)

func (d *defaultFactory) CreateFileFormatter(inStream <-chan []string, outStream chan<- io.ReadCloser) interactor.FileFormatter {
	return NewCsvFormatter(inStream, outStream, d.fileSizeLimit, d.fileLineLimit)
}

func (d *defaultFactory) CreateDataFetcher(outStream chan<- []string) interactor.DataFetcher {
	return &SqlDataFetcher{
		OutStream:   outStream,
		CountQuery:  d.countQuery,
		SelectQuery: d.selectQuery,
		Db:          d.sqlGateway,
	}
}

func (d *defaultFactory) CreateUploader(inStream <-chan io.ReadCloser) interactor.Uploader {
	return NewFtpUploader(d.ftpGateway, d.generationType, inStream)
}

func NewDefaultFactory(
	config FeedConfig,
	sqlGateway SqlGateway,
	ftpGateway FtpGateway,
	generationType string,
) (interactor.FeedFactory, error) {
	countQuery, err := readSqlFromFile(config.CountQueryFilename)
	if err != nil {
		return nil, err
	}
	selectQuery, err := readSqlFromFile(config.SelectQueryFilename)
	if err != nil {
		return nil, err
	}
	return &defaultFactory{
		generationType: generationType,
		fileSizeLimit:  config.FileSizeLimit,
		fileLineLimit:  config.FileLineLimit,
		countQuery:     countQuery,
		selectQuery:    selectQuery,
		sqlGateway:     sqlGateway,
		ftpGateway:     ftpGateway,
	}, nil
}

func readSqlFromFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	sql, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(sql), nil
}

func NewYandexFactory() interactor.FeedFactory {
	return nil
}
