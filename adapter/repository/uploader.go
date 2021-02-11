package repository

import (
	"context"
	"io"

	"github.com/rs/zerolog/log"
)

type (
	FtpGateway interface {
		Upload(ctx context.Context, path string, r io.Reader) error
	}

	ftpUploader struct {
		ftp              FtpGateway
		generationType   string
		inStream         <-chan io.ReadCloser
		uploadedFilesNum uint
		onUpload         func(uploadedFilesNum uint)
	}
)

func (u *ftpUploader) UploadFiles(ctx context.Context) error {
	for {
		select {
		case file, isOpen := <-u.inStream:
			if !isOpen {
				return nil
			}
			if err := u.ftp.Upload(ctx, u.generationType, file); err != nil {
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
