# goAndView

Build android WebView app with Golang:

```bash
#With Android SDK with the NDK bundle installed. Gio currently requires SDK versions >= 31.
#https://developer.android.com/studio#command-tools

#install gofork: https://github.com/gofork-org/go
go install gioui.org/cmd/gogio@latest
export PATH=$PATH:$GOPATH/bin
go get -u github.com/diyism/goAndView
cd $GOPATH/src/github.com/diyism/goAndView/apps/hello
go mod tidy
gogio --target android ./
adb install hello.apk
```
