//+build android

package gowebview

import (
	"git.wow.st/gmp/jni"
	"sync"
	"unsafe"
)

//go:generate javac -source 8 -target 8 -bootclasspath $ANDROID_HOME\platforms\android-29\android.jar -d $TEMP\gowebview\classes gowebview_android.java
//go:generate jar cf gowebview_android.jar -C $TEMP\gowebview\classes .

const (
	// HintNone set the width and height are default size
	HintNone Hint = iota

	// HintFixed prevents the window size to be changed by a user
	HintFixed

	// HintMin set the minimum bounds
	HintMin

	// HintMax set the maximum bounds
	HintMax
)

type webview struct {
	vm   jni.JVM
	view jni.Class

	clsWebView jni.Class
	objWebView jni.Object

	closed chan bool
	mutex  *sync.Mutex
}

func newWindow(config *Config) (wv WebView, err error) {
	w := &webview{
		closed: make(chan bool),
		mutex:  new(sync.Mutex),
	}

	if config.VM == 0 {
		return
	}

	if config.Window == 0 {
		return
	}

	w.vm = jni.JVMFor(config.VM)
	w.view = jni.Class(config.Window)

	err = jni.Do(w.vm, func(env jni.Env) error {
		obj := jni.Object(w.view)
		cls := jni.GetObjectClass(env, obj)

		// Run getClass() to get the Class of w.view (GioView, in case of Gio)
		mid := jni.GetMethodID(env, cls, "getClass", "()Ljava/lang/Class;")
		obj, err := jni.CallObjectMethod(env, obj, mid)
		if err != nil {
			panic(err)
		}

		// Run getClassLoader() to get the ClassLoader
		cls = jni.GetObjectClass(env, obj)
		mid = jni.GetMethodID(env, cls, "getClassLoader", "()Ljava/lang/ClassLoader;")
		obj, err = jni.CallObjectMethod(env, obj, mid)
		if err != nil {
			panic(err)
		}

		// Run findClass() to get the custom class (in that case com.inkeliz.gowebview.gowebview_android, that name
		// is defined on `gowebview_android.java`.
		cls = jni.GetObjectClass(env, obj)
		mid = jni.GetMethodID(env, cls, "findClass", "(Ljava/lang/String;)Ljava/lang/Class;")
		clso, err := jni.CallObjectMethod(env, obj, mid, jni.Value(jni.JavaString(env, `com.inkeliz.gowebview.gowebview_android`)))
		if err != nil {
			panic(err)
		}

		// We need to create an GlobalRef of our class, otherwise we can't manipulate that afterwards.
		w.clsWebView = jni.Class(jni.NewGlobalRef(env, clso))

		// Create a new Object from our class. It's almost the same of `new gowebview_android()`, supposing that we are in
		// Java. The `<init>` and `NewObject` are used to create a "variable" with the class that we get before.
		mid = jni.GetMethodID(env, w.clsWebView, "<init>", `()V`)
		obj, err = jni.NewObject(env, w.clsWebView, mid)
		if err != nil {
			panic(err)
		}

		// We need to create an GlobalRef of our object.
		w.objWebView = jni.Object(jni.NewGlobalRef(env, obj))

		// That is calling the `webview_create` function, which is defined inside the `gowebview_android` class
		// you can view that at gowebview_android.java.
		mid = jni.GetMethodID(env, w.clsWebView, "webview_create", "(Landroid/view/View;)V")
		err = jni.CallVoidMethod(env, w.objWebView, mid, jni.Value(w.view))
		if err != nil {
			panic(err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *webview) Run() {
	w.call("webview_run", "()V")
	<-w.closed
}

func (w *webview) Terminate() {
	return
}

func (w *webview) Destroy() {
	w.call("webview_destroy", "()V")

	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.objWebView == 0 || w.clsWebView == 0 {
		return
	}

	jni.Do(w.vm, func(env jni.Env) error {
		jni.DeleteGlobalRef(env, jni.Object(w.clsWebView))
		jni.DeleteGlobalRef(env, w.objWebView)

		return nil
	})

	w.objWebView, w.clsWebView = 0, 0
	w.closed <- true
}

func (w *webview) Window() uintptr {
	return uintptr(unsafe.Pointer(w.view))
}

func (w *webview) SetTitle(title string) {
	return
}

func (w *webview) SetSize(width int64, height int64, hint Hint) {
	return
}

func (w *webview) SetURL(url string) {
	w.callArgs("webview_navigate", "(Ljava/lang/String;)V", func(env jni.Env) []jni.Value {
		return []jni.Value{
			jni.Value(jni.JavaString(env, url)),
		}
	})
}

func (w *webview) Init(js string) {
	return
}

func (w *webview) Eval(js string) {
	return
}

func (w *webview) call(name, sig string) (err error) {
	// The arguments may need the `env`
	// In that case there's no input, so it's using func(env jni.Env) []jni.Value { return nil } instead
	return w.callArgs(name, sig, func(env jni.Env) []jni.Value { return nil })
}

func (w *webview) callArgs(name, sig string, args func(env jni.Env) []jni.Value) (err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.objWebView == 0 || w.clsWebView == 0 {
		return
	}

	return jni.Do(w.vm, func(env jni.Env) error {
		return jni.CallVoidMethod(env, w.objWebView, jni.GetMethodID(env, w.clsWebView, name, sig), args(env)...)
	})
}

