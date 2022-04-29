# goAndView

Build android WebView app with Golang:

```bash
#With Android SDK with the NDK bundle installed. Gio currently requires SDK versions >= 31.
#https://developer.android.com/studio#command-tools
wget https://dl.google.com/android/repository/commandlinetools-linux-8092744_latest.zip
unzip commandlinetools-linux-8092744_latest.zip
./cmdline-tools/bin/sdkmanager --sdk_root=$ANDROID_SDK_ROOT --install "platforms;android-31"

#install gofork: https://github.com/gofork-org/go
go install gioui.org/cmd/gogio@latest
export PATH=$PATH:$GOPATH/bin
go get -u github.com/diyism/goAndView
cd $GOPATH/src/github.com/diyism/goAndView/apps/hello
go mod tidy
gogio --target android ./
adb install hello.apk
```
