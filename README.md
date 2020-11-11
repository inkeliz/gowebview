# gowebview

[![GoDoc](https://godoc.org/github.com/inkeliz/gowebview?status.svg)](https://godoc.org/github.com/inkeliz/gowebview)
[![Go Report Card](https://goreportcard.com/badge/github.com/zserge/webview)](https://goreportcard.com/report/github.com/inkeliz/gowebview)

A small WebView without CGO, based on [webview/webview](https://github.com/webview/webview). The main goal was to avoid CGO and make it possible to embed the DLLs. Instead of relying directly on CGO as webview/webview, the inkeliz/gowebview uses 
`golang.org/x/sys/windows`, for Windows.

## Why use inkeliz/gowebview?

- If you like to avoid CGO on Windows.
- If you like to have a single `.exe`.

## Why use [webview/webview](https://github.com/webview/webview)?

- If you need support for Darwin/Linux.
- If you need binds from Javascript to Golang.
- If you like to use a more battle-tested and stable library.

### Getting started

Import the package and start using it:

```go
package main

import "github.com/inkeliz/gowebview"

func main() {
	w, err := gowebview.New(&gowebview.Config{URL: "https://google.com", WindowConfig: &gowebview.WindowConfig{Title: "Hello World"}})
	if err != nil {
		panic(err)
	}

	defer w.Destroy()
	w.Run()
}

```

It will open the `https://google.com` webpage, without any additional setup.

## TODO

1. ~~Add support to programmatically allow "locahost connections".~~

    > ~~Currently you must run `CheckNetIsolation.exe LoopbackExempt -a -n="Microsoft.Win32WebViewHost_cw5n1h2txyewy"` externally.~~

    > DONE. It's now implemented for Windows, if you use `Config.TransportConfig{IgnoreNetworkIsolation: true}`. You don't need to execute the sofware as admin, it'll be request only when needed.                                                                                                                                                                                                                                          

2. Improve support for Android (currently it needs [Gio](https://gioui.org))

    > Currently it's supported, but needs Gio to works.

3. Remove dependency of `webview.dll` by calling `WebView2Loader.dll` directly.

    > Currently, it needs to extract both `.dll`.

3. Improve security by restricting where look for DLL.

    > Currently gowebview adds a new filepath on %PATH%)
 
3. Improve error returns.

    > Currently it returns errors directly from `os`, `windows` and `ioutil` libraries.

## License

Code is distributed under MIT license, feel free to use it in your proprietary
projects as well.

## Credits

It's highly based on [webview/webview](https://github.com/webview/webview) and [webview/webview_csharp](https://github.com/webview/webview_csharp). The idea of avoid CGO was inspired by [Ebiten](https://github.com/hajimehoshi/ebiten).

