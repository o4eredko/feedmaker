package repository_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/inhies/go-bytesize"
	"github.com/stretchr/testify/assert"

	"go-feedmaker/adapter/repository"
)

func TestLimitBuffer_Write(t *testing.T) {
	testCases := []struct {
		name      string
		sizeLimit bytesize.ByteSize
		lineLimit uint
		data      []byte
		wantN     int
		wantErr   error
	}{
		{
			name:      "succeed",
			sizeLimit: 11,
			lineLimit: 1,
			data:      []byte("hello world"),
			wantN:     11,
		},
		{
			name:      "size overflow",
			sizeLimit: bytesize.B,
			lineLimit: 100,
			data:      []byte("hello world"),
			wantErr:   repository.ErrSizeOverflow,
		},
		{
			name:      "lines overflow",
			sizeLimit: bytesize.MB,
			lineLimit: 0,
			data:      []byte("hello world"),
			wantErr:   repository.ErrLinesOverflow,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			buffer := new(bytes.Buffer)
			limitWriter := repository.NewLimitWriter(buffer, testCase.sizeLimit, testCase.lineLimit)

			gotN, gotErr := limitWriter.Write(testCase.data)
			assert.NoError(t, limitWriter.Flush())

			assert.Equal(t, testCase.wantN, gotN)
			assert.Equal(t, testCase.wantErr, gotErr)

			if testCase.wantErr == nil {
				got, err := ioutil.ReadAll(buffer)
				assert.NoError(t, err)
				assert.Equal(t, testCase.data, got)
			}
		})
	}
}
