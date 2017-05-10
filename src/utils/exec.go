package utils

import (
	"config"
	"os"
	"os/exec"
)

func Busybox(args ...string) error {
	bb := config.Get("busybox")
	return Run(bb, args...)
}

func Run(cmd string, args ...string) error {
	run := exec.Command(cmd, args...)
	run.Stderr = os.Stderr
	run.Stdout = os.Stdout
	return run.Run()
}
