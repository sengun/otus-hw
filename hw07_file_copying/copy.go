package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	progressBarLength = 50
	copyBufferSize    = 1024
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileIsNotReadable     = errors.New("source file is not readable")
	ErrFileIsNotWritable     = errors.New("destination file is not writable")
	ErrWrongParams           = errors.New("negative offset or limit")
)

func Copy(fromPath, toPath string, offset, limit int64, output io.Writer) error {
	if offset < 0 || limit < 0 {
		return ErrWrongParams
	}

	var sourceFile, destinationFile *os.File

	success := false
	defer func() {
		if !success && destinationFile != nil {
			os.Remove(destinationFile.Name())
		}
	}()

	sourceFile, err := os.Open(fromPath)
	if err != nil {
		return ErrFileIsNotReadable
	}
	defer sourceFile.Close()

	fileInfo, err := sourceFile.Stat()
	if err != nil || !fileInfo.Mode().IsRegular() || fileInfo.Size() == 0 {
		return ErrUnsupportedFile
	}
	fileSize := fileInfo.Size()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	_, err = sourceFile.Seek(offset, 0)
	if err != nil {
		return ErrUnsupportedFile
	}

	destinationFile, err = os.Create(toPath)
	if err != nil {
		return ErrFileIsNotWritable
	}
	defer destinationFile.Close()

	bytesToCopy := fileSize - offset
	if limit == 0 || limit > bytesToCopy {
		limit = bytesToCopy
	}

	reader := io.LimitReader(sourceFile, limit)
	writer := newWriterWithProgressBar(destinationFile, limit, output)
	defer writer.finish()

	buffer := make([]byte, copyBufferSize)
	_, err = io.CopyBuffer(writer, reader, buffer)
	if err != nil {
		return ErrFileIsNotWritable
	}

	success = true
	return nil
}

type WriterWithProgressBar struct {
	destination io.Writer
	output      io.Writer
	byteAmount  int64
	byteWritten int64
}

func (w *WriterWithProgressBar) Write(b []byte) (int, error) {
	n, err := w.destination.Write(b)
	if err != nil {
		return 0, err
	}

	w.byteWritten += int64(n)
	w.printProgress()

	return n, nil
}

func (w *WriterWithProgressBar) printProgress() {
	percent := float32(w.byteWritten) / float32(w.byteAmount)
	filled := strings.Repeat("#", int(progressBarLength*percent))
	empty := strings.Repeat(" ", progressBarLength-len(filled))

	fmt.Fprintf(w.output, "\r[%s%s] %3.1f%%", filled, empty, percent*100)
}

func (w *WriterWithProgressBar) finish() {
	fmt.Fprintln(w.output, "")
}

func newWriterWithProgressBar(file *os.File, totalBytes int64, output io.Writer) *WriterWithProgressBar {
	return &WriterWithProgressBar{
		destination: file,
		output:      output,
		byteAmount:  totalBytes,
	}
}
