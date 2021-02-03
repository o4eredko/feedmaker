package repository_test

import (
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
			buf := repository.NewLimitBuffer(testCase.sizeLimit, testCase.lineLimit)

			gotN, gotErr := buf.Write(testCase.data)

			assert.Equal(t, testCase.wantN, gotN)
			assert.Equal(t, testCase.wantErr, gotErr)
		})
	}
}
