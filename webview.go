package gowebview

//go:generate go run ./generator/generate.go

import (
	"unsafe"
)

// Hints are used to configure window sizing and resizing
type Hint int

type WebView interface {

	// Run runs the main loop until it's terminated. After this function exits -
	// you must destroy the webview.
	Run()

	// Terminate stops the main loop. It is safe to call this function from
	// a background thread.
	Terminate()

	// Destroy destroys a webview and closes the native window.
	Destroy()

	// Window returns a native window handle pointer. When using GTK backend the
	// pointer is GtkWindow pointer, when using Cocoa backend the pointer is
	// NSWindow pointer, when using Win32 backend the pointer is HWND pointer.
	Window() unsafe.Pointer

	// SetTitle updates the title of the native window. Must be called from the UI
	// thread.
	SetTitle(title string)

	// SetSize updates native window size. See Hint constants.
	SetSize(w int64, h int64, hint Hint)

	// Navigate navigates webview to the given URL. URL may be a data URI, i.e.
	// "data:text/text,<html>...</html>". It is often ok not to url-encode it
	// properly, webview will re-encode it for you.
	SetURL(url string)

	// Init injects JavaScript code at the initialization of the new page. Every
	// time the webview will open a the new page - this initialization code will
	// be executed. It is guaranteed that code is executed before window.onload.
	Init(js string)

	// Eval evaluates arbitrary JavaScript code. Evaluation happens asynchronously,
	// also the result of the expression is ignored. Use RPC bindings if you want
	// to receive notifications about the results of the evaluation.
	Eval(js string)
}

type Config struct {
	// Title defines the title of the window
	Title string
	// Size defines the Width x Height of the window
	Size Point
	// Path defines the path where the DLL will be exported
	PathExtraction string
	// IgnoreExtraction if true the DLL will not be exported/extracted. If true you need to keep the DLL in the
	// same folder of the executable or at any path at %PATH%.
	IgnoreExtraction bool
	// Index defines the default page (trigger when SetURL("") and before any SetURL call)
	Index string
	// Debug if is non-zero the Developer Tools will be enabled (if supported)
	Debug bool

	Window uintptr
}

type Point struct {
	X, Y int64
}

// New calls NewWindow to create a new window and a new webview instance. If debug
// is non-zero - developer tools will be enabled (if the platform supports them).
func New(config *Config) (WebView, error) {
	return newWindow(config)
}
