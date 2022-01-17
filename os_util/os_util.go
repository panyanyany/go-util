package os_util

import (
	"os"
	"path/filepath"
)

func MustGetExecutable() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	ex, err = filepath.EvalSymlinks(ex)
	if err != nil {
		panic(err)
	}

	return ex
}

func MustGetExecutableDir() string {
	return filepath.Dir(MustGetExecutable())
}
