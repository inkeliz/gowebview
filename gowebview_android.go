//+build android

package gowebview

import (
	"crypto/x509"
	"encoding/base64"
	"errors"
	"git.wow.st/gmp/jni"
	"sync"
	"unsafe"
)

//to generate gowebview_android.jar:
//cd $GOPATH/src/github.com/diyism/goAndView/
//rm -rf /tmp/gowebview/
//javac -source 8 -target 8 -bootclasspath $ANDROID_SDK_ROOT/platforms/android-29/android.jar -d /tmp/gowebview/classes gowebview_android.java
//jar cf gowebview_android.jar -C /tmp/gowebview/classes .

type webview struct {
	vm   jni.JVM
	view jni.Class

	clsWebView jni.Class
	objWebView jni.Object

	closed chan bool
	mutex  *sync.Mutex

	config *Config
}

func newWindow(config *Config) (wv WebView, err error) {
	w := &webview{
		closed: make(chan bool),
		mutex:  new(sync.Mutex),
		config: config,
	}

	if config.WindowConfig.VM == 0 || config.WindowConfig.Window == 0 {
		return
	}

	w.vm = jni.JVMFor(config.WindowConfig.VM)
	w.view = jni.Class(config.WindowConfig.Window)

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

	w.SetURL(config.URL)
	w.setProxy(config.TransportConfig.Proxy)
	w.setCerts(config.TransportConfig.CertificateAuthorities)

	return w, nil
}

func (w *webview) Hibernate() {
	w.call("webview_hibernate", "()V")
	<-w.closed
}

func (w *webview) Vibrate() {
	w.call("webview_vibrate", "()V")
	<-w.closed
}

func (w *webview) Wakelock() {
	w.call("webview_wakelock", "()V")
	<-w.closed
}

func (w *webview) Locktask() {
	w.call("webview_locktask", "()V")
	<-w.closed
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

func (w *webview) SetSize(point *Point, hint Hint) {
	return
}

func (w *webview) SetURL(url string) {
	if url == "" {
		url = w.config.URL
	}

	w.callArgs("webview_navigate", "(Ljava/lang/String;)V", func(env jni.Env) []jni.Value {
		return []jni.Value{
			jni.Value(jni.JavaString(env, url)),
		}
	})
}

func (w *webview) SetVisibility(v Visibility) {
	switch v {
	case VisibilityMinimized:
		w.call("webview_hide", "()V")
	case VisibilityMaximized:
		w.call("webview_run", "()V")
	default:
		w.call("webview_run", "()V")
	}
}

func (w *webview) setProxy(proxy *HTTPProxy) error {
	if proxy == nil || (proxy.IP == "" && proxy.Port == "") {
		return nil
	}

	ok, err := w.callBooleanArgs("webview_proxy", "(Ljava/lang/String;Ljava/lang/String;)Z", func(env jni.Env) []jni.Value {
		return []jni.Value{
			jni.Value(jni.JavaString(env, proxy.IP)),
			jni.Value(jni.JavaString(env, proxy.Port)),
		}
	})

	if err != nil {
		return err
	}

	if !ok {
		return errors.New("impossible to set proxy")
	}

	return nil
}

func (w *webview) setCerts(certs []x509.Certificate) error {
	if certs == nil {
		return nil
	}

	var jcerts string
	for _, c := range certs {
		jcerts += base64.StdEncoding.EncodeToString(c.Raw) + ";"
	}

	ok, err := w.callBooleanArgs("webview_certs", "(Ljava/lang/String;)Z", func(env jni.Env) []jni.Value {
		return []jni.Value{
			jni.Value(jni.JavaString(env, jcerts)),
		}
	})

	if err != nil {
		return err
	}

	if !ok {
		return errors.New("impossible to set certs")
	}

	return nil
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

func (w *webview) callBooleanArgs(name, sig string, args func(env jni.Env) []jni.Value) (b bool, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.objWebView == 0 || w.clsWebView == 0 {
		return
	}

	err = jni.Do(w.vm, func(env jni.Env) error {
		b, err = jni.CallBooleanMethod(env, w.objWebView, jni.GetMethodID(env, w.clsWebView, name, sig), args(env)...)
		return err
	})

	return b, err
}
