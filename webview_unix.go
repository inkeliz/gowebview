//+build !android,!windows,linux,darwin

package gowebview

/*
#cgo linux openbsd freebsd CXXFLAGS: -DWEBVIEW_GTK -std=c++11
#cgo linux openbsd freebsd pkg-config: gtk+-3.0 webkit2gtk-4.0

#cgo darwin CXXFLAGS: -DWEBVIEW_COCOA -std=c++11
#cgo darwin LDFLAGS: -framework WebKit

#cgo windows CXXFLAGS: -std=c++11

#cgo windows CFLAGS: -g -std=c11 -Iinclude -DPLAIN_API_ONLY -D_BSD_SOURCE -D_DEFAULT_SOURCE

#include "webview_unix.h"

#include <stdlib.h>
#include <stdint.h>
*/
import "C"
import (
	"runtime"
	"sync"
	"unsafe"
)

func init() {
	// Ensure that main.main is called from the main thread
	runtime.LockOSThread()
}

const (
	// HintNone set the width and height are default size
	HintNone = C.WEBVIEW_HINT_NONE

	// HintFixed prevents the window size to be changed by a user
	HintFixed = C.WEBVIEW_HINT_FIXED

	// HintMin set the minimum bounds
	HintMin = C.WEBVIEW_HINT_MIN

	// HintMax set the maximum bounds
	HintMax = C.WEBVIEW_HINT_MAX
)

type webview struct {
	w C.webview_t
}

var (
	m        sync.Mutex
	index    uintptr
	dispatch = map[uintptr]func(){}
	bindings = map[uintptr]func(id, req string) (interface{}, error){}
)

func boolToInt(b bool) C.int {
	if b {
		return 1
	}
	return 0
}

func newWindow(config *Config) (WebView, error) {
	w := new(webview)
	w.w = C.webview_create(boolToInt(config.Debug), config.Window)
	return w, nil
}

func (w *webview) Destroy() {
	C.webview_destroy(w.w)
}

func (w *webview) Run() {
	C.webview_run(w.w)
}

func (w *webview) Terminate() {
	C.webview_terminate(w.w)
}

func (w *webview) Window() uintptr {
	return uintptr(unsafe.Pointer(C.webview_get_window(w.w)))
}

func (w *webview) SetURL(url string) {
	s := C.CString(url)
	defer C.free(unsafe.Pointer(s))
	C.webview_navigate(w.w, s)
}

func (w *webview) SetTitle(title string) {
	s := C.CString(title)
	defer C.free(unsafe.Pointer(s))
	C.webview_set_title(w.w, s)
}

func (w *webview) SetSize(width int64, height int64, hint Hint) {
	C.webview_set_size(w.w, C.int(width), C.int(height), C.int(hint))
}

func (w *webview) Init(js string) {
	s := C.CString(js)
	defer C.free(unsafe.Pointer(s))
	C.webview_init(w.w, s)
}

func (w *webview) Eval(js string) {
	s := C.CString(js)
	defer C.free(unsafe.Pointer(s))
	C.webview_eval(w.w, s)
}
