package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	sourceFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	fileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	totalSize := fileInfo.Size() - offset
	if limit > 0 && limit < totalSize {
		totalSize = limit
	}

	bar := pb.StartNew(int(totalSize))
	defer bar.Finish()

	_, err = sourceFile.Seek(offset, 0)
	if err != nil {
		return err
	}

	destinationFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	barReader := bar.NewProxyReader(sourceFile)

	if limit > 0 {
		_, err = io.CopyN(destinationFile, barReader, limit)
	} else {
		_, err = io.Copy(destinationFile, barReader)
	}

	if err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
	}

	return nil
}
