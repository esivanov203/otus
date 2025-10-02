package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestTempDir(t *testing.T) {
	tmp := os.TempDir()
	wDir := filepath.Join(tmp, "env-dir")
	err := os.MkdirAll(wDir, 0755)
	require.NoError(t, err)
	defer func() {
		err := os.RemoveAll(wDir)
		require.NoError(t, err)
	}()

	createFile := func(path string, content []byte) {
		err := os.WriteFile(path, content, 777)
		require.NoError(t, err)
	}
	files := []struct {
		name        string
		content     string
		expectedVal EnvValue
	}{
		{"BAR", "bar\nВторая строка", EnvValue{"bar", false}},
		{"EMPTY", " \n", EnvValue{"", false}},
		{"HELLO", `"hello"`, EnvValue{`"hello"`, false}},
		{"UNSET", "", EnvValue{"", true}},
		{"BIN", "foo\x00with new line\n", EnvValue{"foo\nwith new line", false}},
		{"TRAILING_SPACES", "value with spaces   \t\n", EnvValue{"value with spaces", false}},
		{"UTF8", "привет\nмир", EnvValue{"привет", false}},
		{"IGNORE=ME", "some value", EnvValue{"", true}},
		{"NULL_START", "\x00second line", EnvValue{"\nsecond line", false}},
	}

	expected := Environment{}

	// ignore list
	err = os.Mkdir(filepath.Join(wDir, "CATALOG"), 777)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(wDir, "IGNORE"), []byte("file"), 0o000)
	require.NoError(t, err)

	//create readable files and expected
	for _, cf := range files {
		createFile(filepath.Join(wDir, cf.name), []byte(cf.content))
		if cf.name != "IGNORE=ME" {
			expected[cf.name] = cf.expectedVal
		}
	}

	env, err := ReadDir(wDir)
	require.NoError(t, err)

	require.Equal(t, expected, env)
}
