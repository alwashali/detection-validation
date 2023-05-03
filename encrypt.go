package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/openpgp"
)

// EncryptFiles traverses the file system starting from the given root folder
// and encrypts all files that match the given filename pattern.
// Each encrypted file will be prepended with "enc_".
// Files that are larger than 2 MB will be skipped.
func EncryptFiles(rootFolder, filenamePattern string) error {
	err := filepath.Walk(rootFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Size() > 2*1024*1024 {
			return nil
		}
		matched, err := filepath.Match(filenamePattern, info.Name())
		if err != nil {
			return err
		}
		if matched {
			inFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer inFile.Close()
			outFile, err := os.Create(filepath.Join(filepath.Dir(path), "enc_"+info.Name()))
			if err != nil {
				return err
			}
			defer outFile.Close()
			w, err := openpgp.SymmetricallyEncrypt(outFile, []byte("password"), nil, nil)
			if err != nil {
				return err
			}
			defer w.Close()
			if _, err = io.Copy(w, inFile); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to encrypt files: %v", err)
	}
	return nil
}
