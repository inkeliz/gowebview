# gowebview

[![GoDoc](https://godoc.org/github.com/inkeliz/gowebview?status.svg)](https://godoc.org/github.com/inkeliz/gowebview)
[![Go Report Card](https://goreportcard.com/badge/github.com/zserge/webview)](https://goreportcard.com/report/github.com/inkeliz/gowebview)

A small WebView without CGO, based on [webview/webview](https://github.com/webview/webview). The main goal is to avoid CGO and make it possible to embed the DLLs. Instead of relying directly on CGO as webview/webview, the inkeliz/gowebview uses 
`golang.org/x/sys/windows` and other natives libraries like `os`. 

## Why use inkeliz/gowebview?

- If you like to avoid CGO.
- If you like to have a single `.exe`.

## Why use [webview/webview](https://github.com/webview/webview)?

- If you need support for Darwin/Linux.
- If you need binds from Javascript to Golang.
- If you like to use a more battle-tested library.

### Getting started

Import the package and start using it:

```go
package main

import "github.com/inkeliz/gowebview"

func main() {
	w, err := gowebview.New(&gowebview.Config{Title: "Hello World", Size: gowebview.Point{X: 800, Y: 800}})
	if err != nil {
		panic(err)
	}
	defer w.Destroy()
	w.SetURL(`https://google.com`)
	w.Run()
}
```

It will open the `https://google.com` webpage, without any additional setup.

### How migrate from [webview/webview](https://github.com/webview/webview) to inkeliz/gowebview?

1. `webview` library has been renamed to `gowebview`.
2. It may not work on Darwin/Linux. Currently, tested only on Windows, other platforms may work based on your luck.
3. It doesn't handle `Bind` (from JS to Golang).
3. `w.Navegate(url string)` has been renamed to `w.SetURL(url string)`.
4. `webview.New(debug bool) webviewer` was modified to `gowebview.New(config *Config) (webviewer, error)`. For instance, to enable debug, you 
must use  `gowebview.New(&Config{Debug: true})` instead. Also, there's a `error` return.

## TODO

1. Add support to programmatically allow "locahost connections".

    > Currently you must run `CheckNetIsolation.exe LoopbackExempt -a -n="Microsoft.Win32WebViewHost_cw5n1h2txyewy"` externally.

2. Add support for Android (compatible with [Gio](https://gioui.org))

    > Currently it's not supported at all.

3. Improve security by restricting where look for DLL.

    > Currently gowebview adds a new filepath on %PATH%)
                                                           >
3. Improve error returns.

    > Currently it returns errors directly from `os`, `windows` and `ioutil` libraries.

## License

Code is distributed under MIT license, feel free to use it in your proprietary
projects as well.

## Credits

It's highly based on [webview/webview](https://github.com/webview/webview) and [webview/webview_csharp](https://github.com/webview/webview_csharp). The idea of avoid CGO was inspired by [Ebiten](https://github.com/hajimehoshi/ebiten).

