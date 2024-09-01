package utils

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

func Cmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func WorkerCount(workerFlag string) int {
	if workerFlag == "default" {
		return runtime.NumCPU()
	}
	count, err := strconv.Atoi(workerFlag)
	if err != nil || count < 0 {
		return 0
	}
	return count
}

func LogFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}
