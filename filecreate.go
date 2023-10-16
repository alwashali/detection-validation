package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func filewrite(fpath string, binPath string) {

	if binPath != "" {
		directory := filepath.Dir(binPath)
		fileName := filepath.Base(binPath)

		copyBinaryTo(directory, fileName)

		args := []string{"createfile", "--path", fpath}
		output := execute(binPath, args)
		fmt.Println(output)
		return
	}

	_, err := os.Stat(fpath)
	if os.IsNotExist(err) {
		dir, file := filepath.Split(fpath)
		fmt.Println(dir, "file:", file)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Println(err)
		}

		file_handle, err := os.Create(fpath)
		file_handle.WriteString("Test file\ncreated by https://github.com/alwashali/Detection-Validation")

		if err != nil {
			log.Println(err)
		}
		file_handle.Close()

	}
}
