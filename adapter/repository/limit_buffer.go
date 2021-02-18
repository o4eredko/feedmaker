package repository

import (
	"bufio"
	"errors"
	"io"

	"github.com/inhies/go-bytesize"
)

var (
	ErrSizeOverflow  = errors.New("limitWriter size overflow")
	ErrLinesOverflow = errors.New("limitWriter lines overflow")
)

type (
	LimitWriter struct {
		lineLimit    uint
		sizeLimit    bytesize.ByteSize
		linesWritten int
		bytesWritten int
		w            *bufio.Writer
	}
)

func NewLimitWriter(w io.Writer, sizeLimit bytesize.ByteSize, lineLimit uint) *LimitWriter {
	return &LimitWriter{
		lineLimit: lineLimit,
		sizeLimit: sizeLimit,
		w:         bufio.NewWriter(w),
	}
}

func (l *LimitWriter) Write(data []byte) (int, error) {
	if l.willOverflowLines() {
		return 0, ErrLinesOverflow
	} else if l.willOverflowSize(len(data)) {
		return 0, ErrSizeOverflow
	}
	if n, err := l.w.Write(data); err != nil {
		return 0, err
	} else {
		l.bytesWritten += n
		l.linesWritten += 1
		return n, nil
	}
}

func (l *LimitWriter) Flush() error {
	return l.w.Flush()
}

func (l *LimitWriter) Reset(w io.Writer) {
	l.Flush()
	l.w.Reset(w)
	l.linesWritten = 0
	l.bytesWritten = 0
}

func (l *LimitWriter) LinesWritten() int {
	return l.linesWritten
}

func (l *LimitWriter) willOverflowLines() bool {
	return l.linesWritten+1 > int(l.lineLimit)
}

func (l *LimitWriter) willOverflowSize(dataLen int) bool {
	return dataLen+l.bytesWritten > int(l.sizeLimit)
}
