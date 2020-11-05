package gowebview

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/inkeliz/gowebview/internal/network"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
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

func TestNewLocalHost(t *testing.T) {
	if err := network.DisablePrivateConnections(); err != nil {
		t.Fatal(err)
	}

	ip := `127.0.0.1:9831`

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Testing Localhost connection"))
	})

	go func(ip string, mux *http.ServeMux){
		if err := http.ListenAndServe(ip, mux); err != nil {
			t.Error(err)
		}
	}(ip, mux)

	time.Sleep(1 * time.Second)

	w, err := New(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Destroy()
	w.SetTitle("Hello World")
	w.SetSize(800, 800, HintMin)
	w.SetURL(`http://`+ip)
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
