package repository

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/rs/zerolog/log"
)

type (
	FtpGateway interface {
		Upload(ctx context.Context, path string, r io.Reader) error
		MakeDir(path string) error
		RemoveDir(dir string) error
		ChangeDir(path string) error
		ChangeDirToParent() error
	}

	ftpUploader struct {
		ftp              FtpGateway
		generationType   string
		inStream         <-chan io.ReadCloser
		uploadedFilesNum uint
		onUpload         func(uploadedFilesNum uint)
	}
)

func NewFtpUploader(ftpGateway FtpGateway, generationType string, inStream <-chan io.ReadCloser) *ftpUploader {
	return &ftpUploader{
		ftp:            ftpGateway,
		generationType: generationType,
		inStream:       inStream,
		onUpload:       func(uploadedFilesNum uint) {},
	}
}

func (u *ftpUploader) UploadFiles(ctx context.Context) error {
	if err := u.ftp.RemoveDir(u.generationType); err != nil {
		log.Error().Err(err).Msgf("Cannot remove dir %s on ftp", u.generationType)
	}
	if err := u.ftp.MakeDir(u.generationType); err != nil {
		return err
	}

	for {
		select {
		case file, isOpen := <-u.inStream:
			if !isOpen {
				return nil
			}

			filename := fmt.Sprintf("%s_%d.csv", u.generationType, u.uploadedFilesNum)
			filename = path.Join(u.generationType, filename)
			if err := u.ftp.Upload(ctx, filename, file); err != nil {
				return err
			}

			u.uploadedFilesNum++
			u.onUpload(u.uploadedFilesNum)
			if err := file.Close(); err != nil {
				log.Error().Err(err).Msgf("Cannot close file after uploading")
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (u *ftpUploader) OnUpload(callback func(uploadedFilesNum uint)) {
	u.onUpload = callback
}
