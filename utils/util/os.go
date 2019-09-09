package util

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir
}

func ExecDir() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		panic(err)
	}
	path, err := filepath.Abs(file)
	if err != nil {
		panic(err)
	}
	splitstr := strings.Split(path, "/")
	return strings.Join(splitstr[:len(splitstr)-1], "/")
}

func Hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hostname
}
