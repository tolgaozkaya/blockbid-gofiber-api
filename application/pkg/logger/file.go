package logger

import (
	"io/fs"
	"os"
	"path"

	"blockchain-smart-tender-platform/pkg/helpers"
)

func file() *os.File {
	var file *os.File
	var perm fs.FileMode
	p := path.Join(helpers.ExePath(), "output.log")
	perm = 0644

	file, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, perm)
	if err != nil {
		Log.Fatalf("Error when creating output Log: %s\n", err.Error())
	}
	return file
}
