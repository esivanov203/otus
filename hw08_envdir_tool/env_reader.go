package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	ret := make(Environment)

	for _, entry := range entries {
		name := entry.Name()
		// файлы с названием, содержащим "=" игнорируем
		if strings.Contains(name, "=") {
			continue
		}
		path := filepath.Join(dir, name)

		ev, err := makeFromFile(path)
		if err == nil {
			ret[name] = ev
		}
	}

	return ret, nil
}

func makeFromFile(path string) (EnvValue, error) {
	f, err := os.Open(path)
	if err != nil {
		return EnvValue{}, err
	}
	defer func() { _ = f.Close() }()

	var fLine string
	val := EnvValue{NeedRemove: true}

	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		fLine = scanner.Text()
		if len(fLine) > 0 {
			val.NeedRemove = false
		}
		fLine = strings.ReplaceAll(fLine, "\x00", "\n")
		fLine = strings.TrimRight(fLine, " \t")
	}
	if errS := scanner.Err(); errS != nil {
		return EnvValue{}, errS
	}

	val.Value = fLine

	return val, nil
}
