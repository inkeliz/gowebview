//+build windows

package gowebview

import (
	"golang.org/x/sys/windows"
	"os"
	"path/filepath"
	"runtime"
	"unsafe"
)

const (
	// Width and height are default size
	HintNone Hint = iota

	// Window size can not be changed by a user
	HintFixed

	// Width and height are minimum bounds
	HintMin

	// Width and height are maximum bounds
	HintMax
)

type webview struct {
	w     uintptr
	dll   *windows.LazyDLL
	index string
}

func newWindow(config *Config) (wv WebView, err error) {
	runtime.LockOSThread()

	if config == nil {
		config = new(Config)
	}

	if config.PathExtraction == "" && config.IgnoreExtraction == false {
		config.PathExtraction = filepath.Join(os.TempDir(), "gowebview")
	}

	if config.IgnoreExtraction == false {
		if err = extract(config.PathExtraction); err != nil {
			return nil, err
		}

		if err = os.Setenv("PATH", config.PathExtraction+`;`+os.Getenv("PATH")); err != nil {
			return nil, err
		}
	}

	w := new(webview)

	w.dll = windows.NewLazyDLL(filepath.Join(config.PathExtraction, "webview.dll"))
	if err = w.dll.Load(); err != nil {
		return nil, err
	}

	w.w, _, err = w.call("webview_create", uintptrBool(config.Debug), config.Window)
	if err != nil && w.w == 0 {
		return nil, err
	}

	if config.Title != "" {
		w.SetTitle(config.Title)
	}

	if config.Index != "" {
		w.index = config.Index
		w.SetURL(config.Index)
	}

	if config.Size.X > 0 && config.Size.Y > 0 {
		w.SetSize(config.Size.X, config.Size.Y, HintMin)
	}

	return w, nil
}

func (w *webview) Run() {
	w.call("webview_run")
}

func (w *webview) Terminate() {
	w.call("webview_terminate")
}

func (w *webview) Destroy() {
	w.call("webview_destroy")
}

func (w *webview) Window() unsafe.Pointer {
	r1, _, _ := w.call("webview_get_window")
	return unsafe.Pointer(&r1)
}

func (w *webview) SetTitle(title string) {
	w.call("webview_set_title", uintptrString(title))
}

func (w *webview) SetSize(width int64, height int64, hint Hint) {
	w.call("webview_set_size", uintptrInt(width), uintptrInt(height), uintptr(hint))
}

func (w *webview) SetURL(url string) {
	if url == "" {
		url = w.index
	}
	w.call("webview_navigate", uintptrString(url))
}

func (w *webview) Init(js string) {
	w.call("webview_init", uintptrString(js))
}

func (w *webview) Eval(js string) {
	w.call("webview_eval", uintptrString(js))
}

func (w *webview) call(function string, a ...uintptr) (uintptr, uintptr, error) {
	if w.w != 0 {
		a = append([]uintptr{w.w}, a...)
	}
	return w.dll.NewProc(function).Call(a...)
}

func uintptrString(s string) uintptr {
	sPtr, _ := windows.BytePtrFromString(s)
	return uintptr(unsafe.Pointer(sPtr))
}

func uintptrInt(i int64) uintptr {
	return uintptr(unsafe.Pointer(&i))
}

func uintptrBool(b bool) uintptr {
	i := 0
	if b {
		i = 1
	}
	return uintptr(unsafe.Pointer(&i))
}
