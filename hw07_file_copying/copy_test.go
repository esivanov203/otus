package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyNegative(t *testing.T) {
	smallF := "/tmp/go_cp_input_small.txt"
	text := "Простой\nкороткий файл"
	f, err := os.Create(smallF)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(smallF)
	}()

	_, err = f.Write([]byte(text))
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)

	errCases := []struct {
		name     string
		expected error
		from     string
		to       string
		offset   int64
		limit    int64
	}{
		{"empty files paths", ErrEmptyFilePath, "", "", 0, 0},
		{"negative offset", ErrNegativeOffsetLimit, "f.in", "f.out", -1, 0},
		{"negative limit", ErrNegativeOffsetLimit, "f.in", "f.out", 0, -1},
		{"input file not exist", ErrFromFileNotFound, "f.in", "f.out", 0, 0},
		{"offset more than input file size", ErrOffsetExceedsFileSize, smallF, "f.out", 100, 0},
		{"unknown input file size", ErrUnsupportedFile, "/dev/urandom", "f.out", 0, 0},
	}

	for _, c := range errCases {
		t.Run(c.name, func(t *testing.T) {
			err := Copy(c.from, c.to, c.offset, c.limit, false)
			require.Equal(t, c.expected, err)
		})
	}

	simpleCases := []struct {
		name     string
		expected string
		from     string
		offset   int64
		limit    int64
	}{
		{"simple", text, smallF, 0, 0},
		{"simple 0 6", text[:6], smallF, 0, 6},
		{"simple 7 0", text[7:], smallF, 7, 0},
		{"simple 7 2", text[7:9], smallF, 7, 2},
		{"simple 0 100 (limit more than size)", text, smallF, 0, 100},
		{"simple 4 100 (limit more than size)", text[4:], smallF, 4, 100},
	}

	for _, c := range simpleCases {
		toPath := "/tmp/test.out"
		t.Run(c.name, func(t *testing.T) {
			err := Copy(c.from, toPath, c.offset, c.limit, false)
			require.NoError(t, err)

			f, err := os.Open(toPath)
			require.NoError(t, err)
			data, err := io.ReadAll(f)
			require.NoError(t, err)
			require.NoError(t, f.Close())

			realS := string(data)
			require.Equal(t, c.expected, realS)

			require.NoError(t, os.Remove(toPath))
		})
	}
}
