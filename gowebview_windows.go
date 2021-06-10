//+build windows,amd64

package gowebview

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/inkeliz/gowebview/internal/network"
	"github.com/inkeliz/gowebview/internal/wincom"
	"github.com/inkeliz/w32"
	"golang.org/x/sys/windows"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

type webview struct {
	dll     *windows.Proc
	browser browser
	view    view
	config  *Config

	done  chan bool
	queue chan func()
}

type browser struct {
	controller *wincom.ICoreWebView2Controller
	webview    *wincom.ICoreWebView2
}

type view struct {
	instance w32.HINSTANCE
	cursor   w32.HCURSOR
	icon     w32.HICON
	window   w32.HWND

	min Point
	max Point
}

func newWindow(config *Config) (wv WebView, err error) {
	w := &webview{
		config: config,
		done:   make(chan bool, 1),
		queue:  make(chan func(), 1<<16),
	}

	if err = extract(config.WindowConfig.Path); err != nil {
		return nil, err
	}

	if config.TransportConfig.IgnoreNetworkIsolation && !network.IsAllowedPrivateConnections() {
		if err := network.EnablePrivateConnections(); err != nil {
			return nil, err
		}
	}

	for _, s := range []string{"WEBVIEW2_BROWSER_EXECUTABLE_FOLDER", "WEBVIEW2_USER_DATA_FOLDER", "WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS", "WEBVIEW2_RELEASE_CHANNEL_PREFERENCE"} {
		os.Unsetenv(s)
	}

	w.setProxy(config.TransportConfig.Proxy)
	w.setCerts(config.TransportConfig.CertificateAuthorities)

	dll, err := windows.LoadDLL(filepath.Join(config.WindowConfig.Path, "WebView2Loader.dll"))
	if err != nil {
		return nil, err
	}

	w.dll, err = dll.FindProc("CreateCoreWebView2EnvironmentWithOptions")
	if err != nil {
		return nil, err
	}

	if err = w.create(); err != nil {
		return nil, err
	}

	w.SetSize(w.config.WindowConfig.Size, HintNone)
	w.SetURL(w.config.URL)
	w.SetTitle(w.config.WindowConfig.Title)

	return w, nil
}

func (w *webview) setProxy(proxy *HTTPProxy) {
	if proxy == nil || (proxy.IP == "" && proxy.Port == "") {
		return
	}

	w.addEnv(` --proxy-server="%s"`, proxy.String())
}

func (w *webview) setCerts(certs []x509.Certificate) {
	if certs == nil || len(certs) == 0 {
		return
	}

	var jcerts string
	h := sha256.New()
	for _, c := range certs {
		h.Write(c.RawSubjectPublicKeyInfo)
		jcerts += base64.StdEncoding.EncodeToString(h.Sum(nil)) + ","
		h.Reset()
	}

	w.addEnv(` --ignore-certificate-errors-spki-list="%s"`, jcerts)
}

func (w *webview) addEnv(argument, value string) {
	os.Setenv("WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS", os.Getenv("WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS")+" "+fmt.Sprintf(argument, value))
}

func (w *webview) create() error {
	cerr := make(chan error, 1<<2)

	go func() {
		runtime.LockOSThread()
		w32.CoInitializeEx(w32.COINIT_APARTMENTTHREADED)

		go func() {
			<-w.done
			cerr <- nil
		}()

		if err := w.createWindow(); err != nil {
			cerr <- err
			return
		}

		watchlist.Store(w.view.window, w)
		defer watchlist.Delete(w.view.window)

		res, _, err := w.dll.Call(0, uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(w.config.WindowConfig.Path))), 0, w.environmentCompletedHandler())
		if res != 0 {
			cerr <- err
			return
		}

		if err := w.loop(); err != nil {
			return
		}
	}()

	return <-cerr
}

func (w *webview) Run() {
	<-w.done
}

func (w *webview) loop() (err error) {
	msg := new(w32.MSG)
loop:
	for {
		select {
		case f := <-w.queue:
			f()
		default:
			m := w32.GetMessage(msg, 0, 0, 0)
			switch m {
			case -1:
				return errors.New("GetMessage fails")
			case 0:
				break loop
			}

			w32.TranslateMessage(msg)
			w32.DispatchMessage(msg)
		}
	}

	return nil
}

func (w *webview) Terminate() {
	w.queue <- func() {
		w32.PostQuitMessage(0)
		w32.DestroyWindow(w.view.window)
		w.done <- true
	}
}

func (w *webview) Destroy() {
	w.Terminate()
}

func (w *webview) Window() uintptr {
	return uintptr(w.view.window)
}

func (w *webview) SetTitle(title string) {
	w.queue <- func() {
		w32.SetWindowText(w.view.window, title)
	}
}

func (w *webview) SetSize(point *Point, hint Hint) {
	defer w.updateSize(false)

	if point == nil {
		return
	}

	switch hint {
	case HintNone:
		w.queue <- func() {
			w32.SetWindowPos(w.view.window, w32.HWND_TOP, 0, 0, int(point.X), int(point.Y), w32.SWP_NOMOVE)
		}
	case HintFixed:
		w.view.min = *point
		w.view.max = *point
	case HintMin:
		w.view.min = *point
	case HintMax:
		w.view.max = *point
	}
}

func (w *webview) updateSize(now bool) {
	if w.browser.webview == nil {
		return
	}

	f := func() { syscall.Syscall(w.browser.controller.VTBL.PutBounds, 2, uintptr(unsafe.Pointer(w.browser.controller)), uintptr(unsafe.Pointer(w32.GetClientRect(w.view.window))), 0) }
	if now {
		f()
		return
	}

	w.queue <- f
}

func (w *webview) SetURL(url string) {
	if url == "" {
		url = w.config.URL
	}

	w.queue <- func() {
		syscall.Syscall(w.browser.webview.VTBL.Navigate, 2, uintptr(unsafe.Pointer(w.browser.webview)), uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(url))), 0)
	}
}

func (w *webview) SetVisibility(v Visibility) {
	switch v {
	case VisibilityMaximized:
		w32.ShowWindow(w.view.window, w32.SW_MAXIMIZE)
	case VisibilityMinimized:
		w32.ShowWindow(w.view.window, w32.SW_MINIMIZE)
	default:
		w32.ShowWindow(w.view.window, w32.SW_SHOWDEFAULT)
		w32.SetForegroundWindow(w.view.window)
	}
}

// watchlist is kinda of `map[hwnd]*webview
var watchlist sync.Map

func watch(hwnd w32.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	ww, ok := watchlist.Load(hwnd)
	if !ok {
		return w32.DefWindowProc(hwnd, msg, wParam, lParam)
	}

	w := ww.(*webview)

	switch msg {
	case w32.WM_SIZE, w32.WM_SIZING, w32.WM_WINDOWPOSCHANGED:
		w.updateSize(true)
	case w32.WM_ERASEBKGND:
		return 1
	case w32.WM_DESTROY, w32.WM_CLOSE:
		w.Destroy()
	case w32.WM_PAINT:
		p := new(w32.PAINTSTRUCT)
		w32.BeginPaint(hwnd, p)
		w32.EndPaint(hwnd, p)
	case w32.WM_GETMINMAXINFO:
		mm := (*w32.MINMAXINFO)(unsafe.Pointer(lParam))
		if w.view.min.X > 0 || w.view.min.Y > 0 {
			mm.PtMinTrackSize = w32.POINT{
				X: int32(w.view.min.X),
				Y: int32(w.view.min.Y),
			}
		}
		if w.view.max.X > 0 || w.view.max.Y > 0 {
			mm.PtMaxTrackSize = w32.POINT{
				X: int32(w.view.max.X),
				Y: int32(w.view.max.Y),
			}
		}
	}

	return w32.DefWindowProc(hwnd, msg, wParam, lParam)
}

func (w *webview) createWindow() error {
	w.view.instance = w32.GetModuleHandle("")
	if w.view.instance == 0 {
		return errors.New("GetModuleHandle fails")
	}

	w.view.cursor = w32.LoadCursorInt(w32.IDC_ARROW)
	if w.view.cursor == 0 {
		return errors.New("LoadCursorInt fails")
	}

	if path, err := os.Executable(); err == nil {
		w.view.icon = w32.ExtractIcon(path, 0)
	}

	if _, ok := w32.GetClassInfoEx(w.view.instance, "webview"); !ok {
		class := w32.RegisterClassEx(&w32.WNDCLASSEX{
			Style:      w32.CS_HREDRAW | w32.CS_VREDRAW | w32.CS_OWNDC,
			WndProc:    windows.NewCallback(watch),
			Instance:   w.view.instance,
			Cursor:     w.view.cursor,
			Icon:       w.view.icon,
			IconSm:     w.view.icon,
			ClassName:  windows.StringToUTF16Ptr("webview"),
			Background: w32.WHITE_BRUSH,
		})
		if class == 0 {
			return errors.New("RegisterClassEx fails")
		}
	}

	w.view.window = w32.CreateWindowEx(
		w32.CS_HREDRAW|w32.CS_VREDRAW|w32.CS_OWNDC,
		windows.StringToUTF16Ptr("webview"),
		windows.StringToUTF16Ptr(""),
		uint(w32.WS_OVERLAPPEDWINDOW),
		w32.CW_USEDEFAULT, w32.CW_USEDEFAULT,
		int(w.config.WindowConfig.Size.X), int(w.config.WindowConfig.Size.Y),
		0,
		0,
		w.view.instance,
		nil,
	)
	if w.view.window == 0 {
		return errors.New("CreateWindowEx failed")
	}

	w.SetVisibility(w.config.WindowConfig.Visibility)
	w32.SetForegroundWindow(w.view.window)
	w32.SetFocus(w.view.window)
	w32.UpdateWindow(w.view.window)

	return nil
}

func (w *webview) controllerCompletedHandler() uintptr {
	h := &wincom.ICoreWebView2CreateCoreWebView2ControllerCompletedHandler{
		VTBL: &wincom.ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerVTBL{
			Invoke: windows.NewCallback(func(i *wincom.ICoreWebView2CreateCoreWebView2ControllerCompletedHandler, p uintptr, createdController *wincom.ICoreWebView2Controller) uintptr {

				syscall.Syscall(createdController.VTBL.AddRef, 1, uintptr(unsafe.Pointer(createdController)), 0, 0)
				w.browser.controller = createdController

				createdWebView2 := new(wincom.ICoreWebView2)

				syscall.Syscall(createdController.VTBL.GetCoreWebView2, 2, uintptr(unsafe.Pointer(createdController)), uintptr(unsafe.Pointer(&createdWebView2)), 0)
				w.browser.webview = createdWebView2

				syscall.Syscall(w.browser.webview.VTBL.AddRef, 1, uintptr(unsafe.Pointer(w.browser.webview)), 0, 0)
				w.done <- true

				return 0
			}),
		},
	}
	h.VTBL.BasicVTBL = wincom.NewBasicVTBL(&h.Basic)
	return uintptr(unsafe.Pointer(h))
}

func (w *webview) environmentCompletedHandler() uintptr {
	h := &wincom.ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler{
		VTBL: &wincom.ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVTBL{
			Invoke: windows.NewCallback(func(i uintptr, p uintptr, createdEnvironment *wincom.ICoreWebView2Environment) uintptr {
				syscall.Syscall(createdEnvironment.VTBL.CreateCoreWebView2Controller, 3, uintptr(unsafe.Pointer(createdEnvironment)), uintptr(w.view.window), w.controllerCompletedHandler())
				return 0
			}),
		},
	}
	h.VTBL.BasicVTBL = wincom.NewBasicVTBL(&h.Basic)
	return uintptr(unsafe.Pointer(h))
}
