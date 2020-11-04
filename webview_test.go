package gowebview

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	w, err := New(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Destroy()
	w.SetTitle("Hello World")
	w.SetSize(800, 800, HintMin)
	w.SetURL(`https://google.com`)
	w.Run()
}

func TestNewConfig(t *testing.T) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(os.TempDir(), hex.EncodeToString(b))

	w, err := New(&Config{Title: "Hello World", Size: Point{X: 800, Y: 800}, PathExtraction: path})
	if err != nil {
		t.Fatal(err)
	}
	defer func(w WebView) {
		w.Destroy()
	}(w)
	w.SetURL(`https://google.com`)
	w.Run()
}
