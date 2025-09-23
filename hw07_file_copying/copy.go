package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/sys/unix"
)

var (
	ErrEmptyFilePath         = errors.New("'from' or 'to' filepath is empty")
	ErrNegativeOffsetLimit   = errors.New("'offset' and 'limit' can't be negative")
	ErrFromFileNotFound      = errors.New("'from' file not found")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrCreateNewFile         = errors.New("create output file")
	ErrReadFile              = errors.New("read input file")
	ErrWriteFile             = errors.New("write output file")
)

func validate(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrEmptyFilePath
	}
	if offset < 0 || limit < 0 {
		return ErrNegativeOffsetLimit
	}
	return nil
}

func loadFile(fromPath string, offset int64) (*os.File, int64, error) {
	fIn, err := os.Open(fromPath)
	if os.IsNotExist(err) {
		return nil, 0, ErrFromFileNotFound
	}
	if err != nil {
		return nil, 0, err
	}

	fInStat, err := fIn.Stat()
	if err != nil {
		return nil, 0, err
	}
	if !fInStat.Mode().IsRegular() {
		return nil, 0, ErrUnsupportedFile
	}
	fInSize := fInStat.Size()
	if offset > fInSize {
		return nil, 0, ErrOffsetExceedsFileSize
	}

	return fIn, fInSize, nil
}

func Copy(fromPath, toPath string, offset, limit int64, showProgress bool) error {
	if err := validate(fromPath, toPath, offset, limit); err != nil {
		return err
	}

	fIn, fInSize, err := loadFile(fromPath, offset)
	if err != nil {
		return err
	}
	defer func() {
		if e := fIn.Close(); e != nil {
			log.Printf("Error closing input file: %v", e)
		}
	}()

	fOut, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCreateNewFile, err)
	}

	defer func() {
		if e := fOut.Close(); e != nil {
			log.Fatalf("error closing out file: %v", e)
		}
	}()

	if offset > 0 {
		fd := int(fIn.Fd())
		_, err := unix.Seek(fd, offset, 1)
		if err != nil {
			return err
		}
	}

	total := fInSize - offset
	if limit != 0 && total > limit {
		total = limit
	}
	var counter int64

	var bufSize int64 = 1024 // для больших файлов лучше 32 * 1024, возможно тоже получать через flag
	if bufSize > total {
		bufSize = total
	}
	buf := make([]byte, bufSize)

	for counter < total {
		if bufSize > total-counter {
			buf = buf[:total-counter]
		}
		readN, errR := fIn.Read(buf)
		if errR != nil && !errors.Is(errR, io.EOF) {
			return fmt.Errorf("%w: %w", ErrReadFile, errR)
		}

		wroteN, errW := fOut.Write(buf[:readN])
		if errW != nil {
			return fmt.Errorf("%w: %w", ErrWriteFile, errW)
		}
		if readN != wroteN {
			return ErrWriteFile
		}
		counter += int64(wroteN)
		if showProgress {
			fmt.Printf("\rProcessed...%0.f%%", float64(counter)/float64(total)*100)
		}
		if errors.Is(errR, io.EOF) {
			break
		}
	}

	return nil
}
