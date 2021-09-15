# goAndView

Build android WebView app with Golang:

```bash
#install forkgo: https://github.com/forkgo-org/go
go install gioui.org/cmd/gogio
export PATH=$PATH:$GOPATH/bin
go get -u github.com/diyism/goAndView
cd $GOPATH/src/github.com/diyism/goAndView/apps/hello
go mod tidy
gogio --target android ./
adb install hello.apk
```
