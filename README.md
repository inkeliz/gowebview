# goAndView

Build android WebView app with Golang:

```bash
#With Android SDK with the NDK bundle installed. Gio currently requires SDK versions >= 31.
#https://developer.android.com/studio#command-tools
wget https://dl.google.com/android/repository/commandlinetools-linux-8092744_latest.zip
unzip commandlinetools-linux-8092744_latest.zip
./cmdline-tools/bin/sdkmanager --sdk_root=$ANDROID_SDK_ROOT --install "platforms;android-31"

#install gofork: https://github.com/gofork-org/goFork
go install gioui.org/cmd/gogio@latest
export PATH=$PATH:$GOPATH/bin
go get -u github.com/diyism/goAndView
cd $GOPATH/src/github.com/diyism/goAndView/apps/hello
go mod tidy

(if you want to add an android permission, for example "android.permission.VIBRATE", "android.permission.WAKE_LOCK",
don't use 'import (_ "gioui.org/app/permission/wakelock")', "go mod tidy" will replace gioui.org v0.0.0-20220414170908-ad7c1eb with v0.3.0, goAndView building will fail.
you can add it into the following file, just after "android.permission.INTERNET":
$GOPATH/pkg/mod/gioui.org/cmd@v0.0.0-20230822165948-7cb98d0557e7/gogio/permission.go
then "cd $GOPATH/pkg/mod/gioui.org/cmd@v0.0.0-20230822165948-7cb98d0557e7/gogio/" and run "go build .",
then "cp gogio $GOPATH/bin/"
)

gogio --target android ./
//gogio --target android --arch arm64 ./         //the hello.apk is only 5.5MB, if only for your arm64 phone
adb install hello.apk
```
