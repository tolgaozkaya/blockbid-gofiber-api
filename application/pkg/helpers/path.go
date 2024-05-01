package helpers

import (
	"os"
	"path/filepath"
)

func ExePath() string {
	return filepath.Dir(Executable())
}

func Executable() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return ex
}
