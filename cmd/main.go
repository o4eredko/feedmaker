package main

import (
	"context"
	"io"
	"time"

	"github.com/inhies/go-bytesize"

	"go-feedmaker/adapter/repository"
)

func main() {
	inStream := make(chan []string)
	outStream := make(chan io.ReadCloser)
	buffer := repository.NewLimitBuffer(bytesize.B*27, 3)
	formatter := repository.NewCsvFormatter(inStream, outStream, buffer)
	go func() {
		if err := formatter.FormatFiles(context.Background()); err != nil {
			panic(err)
		}
	}()
	go func() {
		inStream <- []string{"a1", "b1", "c1"}
		time.Sleep(time.Millisecond * 100)

		inStream <- []string{"a2", "b2", "c2"}
		time.Sleep(time.Millisecond * 100)

		inStream <- []string{"a3", "b3", "c3"}
		time.Sleep(time.Millisecond * 100)

		inStream <- []string{"a4", "b4", "c4"}
		time.Sleep(time.Millisecond * 100)

		inStream <- []string{"a5", "b5", "c5"}
		time.Sleep(time.Millisecond * 100)

		inStream <- []string{"a6", "b6", "c6"}
		time.Sleep(time.Millisecond * 100)

		close(inStream)
	}()
	for file := range outStream {
		defer file.Close()
	}
}
