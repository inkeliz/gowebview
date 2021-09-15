package gowebview

import (
	"crypto/x509"
	"os"
	"path/filepath"
	"strings"
)

//go:generate go run ./generator/generate.go

// WebView is the interface implemented by each webview.
type WebView interface {

	// Run runs the main loop until it's terminated. After this function exits -
	// you must destroy the webview.
	Run()

	Hibernate()

	// Terminate stops the main loop. It is safe to call this function from
	// a background thread.
	Terminate()

	// Destroy destroys a webview and closes the native window.
	Destroy()

	// Window returns a native window handle pointer. When using GTK backend the
	// pointer is GtkWindow pointer, when using Cocoa backend the pointer is
	// NSWindow pointer, when using Win32 backend the pointer is HWND pointer.
	Window() uintptr

	// SetTitle updates the title of the native window. Must be called from the UI
	// thread.
	SetTitle(title string)

	// SetSize updates native window size. See Hint constants.
	SetSize(point *Point, hint Hint)

	// Navigate navigates webview to the given URL. URL may be a data URI, i.e.
	// "data:text/text,<html>...</html>".
	SetURL(url string)

	// SetVisibility updates the WindowMode, such as minimized or maximized
	SetVisibility(v Visibility)

	// Init injects JavaScript code at the initialization of the new page. Every
	// time the webview will open a the new page - this initialization code will
	// be executed. It is guaranteed that code is executed before window.onload.
	//Init(js string)

	// Eval evaluates arbitrary JavaScript code. Evaluation happens asynchronously,
	// also the result of the expression is ignored. Use RPC bindings if you want
	// to receive notifications about the results of the evaluation.
	//Eval(js string)
}

// New calls NewWindow to create a new window and a new webview instance. If debug
// is non-zero - developer tools will be enabled (if the platform supports them).
func New(config *Config) (WebView, error) {
	if config == nil {
		config = new(Config)
	}

	if config.WindowConfig == nil {
		config.WindowConfig = &WindowConfig{}
	}

	if config.TransportConfig == nil {
		config.TransportConfig = &TransportConfig{}
	}

	if config.WindowConfig.Title == "" {
		dir, err := os.Executable()
		if err != nil {
			dir = "gowebview"
		}
		filename := filepath.Base(filepath.Clean(dir))
		config.WindowConfig.Title = strings.Title(strings.TrimSuffix(filename, filepath.Ext(filename)))
	}

	if config.WindowConfig.Size == nil {
		config.WindowConfig.Size = &Point{X: 600, Y: 600}
	}

	if config.WindowConfig.Path == "" {
		config.WindowConfig.Path = filepath.Join(os.TempDir(), config.WindowConfig.Title)
	}

	return newWindow(config)
}


// Config are used to set the initial and default values to the WebView.
type Config struct {

	// WindowConfig keeps configurations about the window
	WindowConfig *WindowConfig

	// TransportConfig keeps configurations about the network traffic
	TransportConfig *TransportConfig

	// URL defines the default page.
	URL string

	// Debug if is non-zero the Developer Tools will be enabled (if supported).
	Debug bool
}

// WindowConfig describes topics related to the Window/View.
type WindowConfig struct {

	// Title defines the title of the window.
	Title string

	// Size defines the Width x Height of the window.
	Size *Point

	// Path defines the path where the DLL will be exported.
	Path string

	// Visibility defines how the page must open.
	Visibility Visibility

	// Window defines the window handle (GtkWindow, NSWindow, HWND pointer or View pointer for Android).
	// For Gio (Android):  it MUST point to `e.View` from `app.ViewEvent`
	Window uintptr

	// VM defines the JNI VM for Android
	// For Gio (Android):  it MUST point to `app.JavaVM()`
	VM uintptr
}

// WindowConfig describes topics related to the network traffic.
type TransportConfig struct {

	// Proxy defines the proxy which the connection will be pass through.
	Proxy *HTTPProxy

	// InsecureBypassCustomProxy if true it might ignore the Proxy settings provided, if fails it to be used. Otherwise
	// it could return ErrFeatureNotSupported or ErrImpossibleProxy when create the WebView.
	//
	// It doesn't have any effect if Proxy is undefined (nil).
	// WARNING: It's might be danger if you define an custom Proxy, since it expose the connection without any proxy.
	InsecureIgnoreCustomProxy bool

	// Authorities defines custom authorities that your app trusts. It could be combined with standard authorities from
	// the machine, it might adds the certificate persistent on the user machine and might requires user interaction
	// to approve such action.
	CertificateAuthorities []x509.Certificate

	// @TODO Support CertificatePinning
	// CertificatePinning if set will drop the connection if connect to one website that doesn't have the exactly
	// same certificate hash.
	// CertificateKeyPinning []byte

	// @TODO Support InsecureIgnoreCAVerification
	// InsecureIgnoreCAVerification if true will load pages can be loaded without verifies the certificate.
	// WARNING: It's might be danger and expose you to MITM.
	//InsecureIgnoreCertificateVerification bool

	// IgnoreNetworkIsolation if true will load pages from the local-network (such as 127.0.0.1).
	IgnoreNetworkIsolation bool
}

// Hint are used to configure window sizing and resizing
type Hint int

const (
	// HintNone set the current and default width and height
	HintNone Hint = iota

	// HintFixed prevents the window size to be changed by a user
	HintFixed

	// HintMin set the minimum bounds
	HintMin

	// HintMax set the maximum bounds
	HintMax
)

// Visibility are used to configure if the window mode (maximized or minimized)
type Visibility int

const (
	// VisibilityDefault will open the window at their default WindowMode
	VisibilityDefault Visibility = iota

	// VisibilityMaximized will open the window as maximized
	// Windowless systems (like Android) will open as fullscreen
	VisibilityMaximized

	// VisibilityMinimized will open the window as minimized
	// Windowless systems (like Android) will hides the webview (return the previous view)
	VisibilityMinimized
)

// Point are used to configure the size or coordinates.
type Point struct {
	X, Y int64
}

// HTTPProxy are used to configure the Proxy
type HTTPProxy struct {
	IP   string
	Port string
}

// Network implements net.Addr
func (p *HTTPProxy) Network() string {
	return "tcp"
}

// String implements net.Addr
func (p *HTTPProxy) String() string {
	if p == nil || (p.Port == "" && p.IP == "") {
		return ""
	}

	if strings.Index(p.IP, `:`) >= 0 && !strings.HasPrefix(p.IP, `[`) && !strings.HasPrefix(p.IP, `]`) {
		return "[" + p.IP + "]:" + p.Port
	}
	return p.IP + ":" + p.Port
}
