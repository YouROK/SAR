package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CheckSu() {
	if os.Getuid() != 0 {
		fmt.Println("Must be run under a root")
		os.Exit(1)
	}
	fmt.Println("Root access granted")
}

func CheckBusyBox() string {
	buf, err := exec.Command("which", "busybox").Output()
	if err != nil {
		fmt.Println("Busybox not found:", err)
		os.Exit(1)
	}
	if len(buf) == 0 {
		fmt.Println("Busybox not found")
		os.Exit(1)
	}
	return strings.TrimSpace(string(buf))
}

func Remount() {
	exec.Command("mount", "-o", "rw,remount", "/system").Run()
}
