package main

import (
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExitCodes(t *testing.T) {
	cases := []struct {
		name         string
		command      []string
		expectedCode int
	}{
		{"normal", []string{"ls"}, 0},
		{"nothing command", []string{"false"}, 1},
		{"ls with invalid path", []string{"ls", "/path/does/not/exist"}, 2},
		{"empty command", []string{}, -1},
		{"handle script with exit 5", []string{"sh", "-c", "exit 5"}, 5},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := RunCmd(c.command, Environment{})
			require.Equal(t, c.expectedCode, actual)
		})
	}
}

func TestSplitEnv(t *testing.T) {
	err := os.Setenv("BAR", "bar")
	require.NoError(t, err)
	err = os.Setenv("BAZ", "baz")
	require.NoError(t, err)

	env := os.Environ()
	countS := len(env)

	fEnv := Environment{
		"BAR": EnvValue{"new_bar", false},
		"BAZ": EnvValue{"new_baz", true},
		"FOO": EnvValue{"new_foo", false},
		"FOZ": EnvValue{"new_foz", true},
	}

	cEnv := splitEnv(fEnv)
	countE := len(cEnv)

	// ожидаем, что в результирующий env - количество не изменилось:
	//   одна запись добавилась,
	//   одна удалилась
	//   и у одной изменилось значение

	require.Equal(t, countS, countE)
	require.True(t, slices.Contains(cEnv, "FOO=new_foo"))
	require.False(t, slices.Contains(cEnv, "BAZ=new_baz"))
	require.True(t, slices.Contains(cEnv, "BAR=new_bar"))
}
