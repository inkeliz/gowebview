// +build windows

package gowebview

import (
	"crypto/subtle"
	"golang.org/x/crypto/blake2b"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// extract will save all DLLs (or equivalent) into the given path. If the path is `C:\Something`, it will create exactly
// C:\Something\WebView2Loader.dll.
//
// The file is extracted by default when use `New(nil)`.
func extract(path string) error {
	var approved int

	blake, _ := blake2b.New256(nil)
	for n, h := range FilesHashes {
		blake.Reset()

		file, err := os.Open(filepath.Join(path, n))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return err
		}

		_, err = io.Copy(blake, file)
		if err != nil {
			return err
		}

		if subtle.ConstantTimeCompare(h, blake.Sum(nil)) == 1 {
			approved += 1
		}
	}

	if approved == len(Files) {
		return nil
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	for n, f := range Files {
		if _, err := os.Stat(filepath.Join(path, n)); !os.IsNotExist(err) {
			continue
		}

		if err := ioutil.WriteFile(filepath.Join(path, n), f, 755); err != nil {
			return err
		}
	}

	return nil
}
