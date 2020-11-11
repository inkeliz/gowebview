//+build windows,amd64

package gowebview

import (
	"github.com/inkeliz/gowebview/internal/network"
	"golang.org/x/sys/windows"
	"os"
	"path/filepath"
	"runtime"
	"unsafe"
)

const (
	// HintNone set the width and height are default size
	HintNone Hint = iota

	// HintFixed prevents the window size to be changed by a user
	HintFixed

	// HintMin set the minimum bounds
	HintMin

	// HintMax set the maximum bounds
	HintMax
)

type webview struct {
	w   uintptr
	dll *windows.LazyDLL

	config *Config
}

func newWindow(config *Config) (wv WebView, err error) {
	runtime.LockOSThread()

	w := &webview{
		config: config,
	}

	if err = extract(config.WindowConfig.Path); err != nil {
		return nil, err
	}

	if err = os.Setenv("PATH", config.WindowConfig.Path+`;`+os.Getenv("PATH")); err != nil {
		return nil, err
	}

	w.dll = windows.NewLazyDLL(filepath.Join(config.WindowConfig.Path, "webview.dll"))
	if err = w.dll.Load(); err != nil {
		return nil, err
	}

	w.w, _, err = w.call("webview_create", uintptrBool(config.Debug), config.WindowConfig.Window)
	if err != nil && w.w == 0 {
		return nil, err
	}

	w.SetTitle(config.WindowConfig.Title)
	w.SetURL(config.URL)
	w.SetSize(config.WindowConfig.Size, HintMin)

	if config.TransportConfig.IgnoreNetworkIsolation && !network.IsAllowedPrivateConnections() {
		if err := network.EnablePrivateConnections(); err != nil {
			return nil, err
		}
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

func (w *webview) Window() uintptr {
	r1, _, _ := w.call("webview_get_window")
	return r1
}

func (w *webview) SetTitle(title string) {
	w.call("webview_set_title", uintptrString(title))
}

func (w *webview) SetSize(point *Point, hint Hint) {
	w.call("webview_set_size", uintptrInt(point.X), uintptrInt(point.Y), uintptr(hint))
}

func (w *webview) SetURL(url string) {
	if url == "" {
		url = w.config.URL
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
