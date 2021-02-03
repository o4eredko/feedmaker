package repository

import (
	"bytes"
	"errors"

	"github.com/inhies/go-bytesize"
)

var (
	ErrSizeOverflow  = errors.New("buffer size overflow")
	ErrLinesOverflow = errors.New("buffer lines overflow")
)

type (
	LimitBuffer struct {
		lineLimit    uint
		sizeLimit    int
		linesWritten int
		buf          *bytes.Buffer
	}
)

func NewLimitBuffer(sizeLimit bytesize.ByteSize, lineLimit uint) *LimitBuffer {
	return &LimitBuffer{
		lineLimit: lineLimit,
		sizeLimit: int(sizeLimit),
		buf:       bytes.NewBuffer(make([]byte, 0, sizeLimit)),
	}
}

func (l *LimitBuffer) Read(p []byte) (n int, err error) {
	return l.buf.Read(p)
}

func (l *LimitBuffer) Write(data []byte) (int, error) {
	if l.willOverflowLines() {
		return 0, ErrLinesOverflow
	} else if l.willOverflowSize(len(data)) {
		return 0, ErrSizeOverflow
	}
	if n, err := l.buf.Write(data); err != nil {
		return 0, err
	} else {
		l.linesWritten += 1
		return n, nil
	}
}

func (l *LimitBuffer) Reset() {
	l.buf.Reset()
	l.linesWritten = 0
}

func (l *LimitBuffer) LinesWritten() int {
	return l.linesWritten
}

func (l *LimitBuffer) willOverflowLines() bool {
	return l.linesWritten+1 > int(l.lineLimit)
}

func (l *LimitBuffer) willOverflowSize(dataLen int) bool {
	return dataLen+l.buf.Len() > l.sizeLimit
}
