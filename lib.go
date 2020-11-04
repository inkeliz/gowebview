package gowebview

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// extract will save all DLLs (or equivalent) into the given path. If the path is `C:\Something`, it will create exactly
// C:\Something\webview.dll and C:\Something\WebView2Loader.dll.
//
// The file is extracted by default when use `New(nil)`.
func extract(path string) error {

	// @TODO Verify hashes of the DLL
	// It need something to enforce and restrict the load of the DLL to one specific path (current it searchers multiple
	// folders).
	/*
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
	*/

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
