package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 || strings.TrimSpace(cmd[0]) == "" {
		return -1
	}

	c := cmd[0]
	params := cmd[1:]
	command := exec.Command(c, params...)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	command.Env = splitEnv(env)

	err := command.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		fmt.Println(err)
		return 1
	}

	return 0
}

func splitEnv(env Environment) []string {
	var cEnv []string

	if env != nil {
		for _, kv := range os.Environ() {
			parts := strings.SplitN(kv, "=", 2) // Разделяем на ключ и значение только один раз
			k := parts[0]
			v := ""
			if len(parts) > 1 {
				v = parts[1]
			}
			if _, ok := env[k]; !ok {
				cEnv = append(cEnv, fmt.Sprintf("%s=%s", k, v))
			}
		}
		for nk, nv := range env {
			if !nv.NeedRemove {
				cEnv = append(cEnv, fmt.Sprintf("%s=%s", nk, nv.Value))
			}
		}
	}

	return cEnv
}
