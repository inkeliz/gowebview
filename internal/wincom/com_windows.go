package wincom

import (
	"golang.org/x/sys/windows"
)

type (
	// Basic is the basic struct that any other must include, which implements basic COM functions and implements NewHandler
	Basic struct{}

	// BasicVTBL is the basic VTBL, which implements basic COM functions, which implements VTBL
	BasicVTBL struct {
		QueryInterface uintptr
		AddRef         uintptr
		Release        uintptr
	}
)

func NewBasicVTBL(h *Basic) BasicVTBL {
	return BasicVTBL{
		QueryInterface: windows.NewCallback(h.QueryInterface),
		AddRef:         windows.NewCallback(h.AddRef),
		Release:        windows.NewCallback(h.Release),
	}
}

// QueryInterface is the QueryInterface from COM
func (obj *Basic) QueryInterface(_ *ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, _, _ uintptr) uintptr {
	return 0
}

// AddRef is the AddRef from COM
func (obj *Basic) AddRef(_ *ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) uintptr {
	return 1
}

// Release is the Release from COM
func (obj *Basic) Release(_ *ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) uintptr {
	return 1
}

type (
	// ICoreWebView2 implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2?view=webview2-1.0.622.22
	ICoreWebView2 struct {
		Basic
		VTBL *ICoreWebView2VTBL
	}

	// ICoreWebView2VTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2?view=webview2-1.0.622.22
	ICoreWebView2VTBL struct {
		BasicVTBL
		GetSettings                            uintptr
		GetSource                              uintptr
		Navigate                               uintptr
		NavigateToString                       uintptr
		AddNavigationStarting                  uintptr
		RemoveNavigationStarting               uintptr
		AddContentLoading                      uintptr
		RemoveContentLoading                   uintptr
		AddSourceChanged                       uintptr
		RemoveSourceChanged                    uintptr
		AddHistoryChanged                      uintptr
		RemoveHistoryChanged                   uintptr
		AddNavigationCompleted                 uintptr
		RemoveNavigationCompleted              uintptr
		AddFrameNavigationStarting             uintptr
		RemoveFrameNavigationStarting          uintptr
		AddFrameNavigationCompleted            uintptr
		RemoveFrameNavigationCompleted         uintptr
		AddScriptDialogOpening                 uintptr
		RemoveScriptDialogOpening              uintptr
		AddPermissionRequested                 uintptr
		RemovePermissionRequested              uintptr
		AddProcessFailed                       uintptr
		RemoveProcessFailed                    uintptr
		AddScriptToExecuteOnDocumentCreated    uintptr
		RemoveScriptToExecuteOnDocumentCreated uintptr
		ExecuteScript                          uintptr
		CapturePreview                         uintptr
		Reload                                 uintptr
		PostWebMessageAsJSON                   uintptr
		PostWebMessageAsString                 uintptr
		AddWebMessageReceived                  uintptr
		RemoveWebMessageReceived               uintptr
		CallDevToolsProtocolMethod             uintptr
		GetBrowserProcessID                    uintptr
		GetCanGoBack                           uintptr
		GetCanGoForward                        uintptr
		GoBack                                 uintptr
		GoForward                              uintptr
		GetDevToolsProtocolEventReceiver       uintptr
		Stop                                   uintptr
		AddNewWindowRequested                  uintptr
		RemoveNewWindowRequested               uintptr
		AddDocumentTitleChanged                uintptr
		RemoveDocumentTitleChanged             uintptr
		GetDocumentTitle                       uintptr
		AddHostObjectToScript                  uintptr
		RemoveHostObjectFromScript             uintptr
		OpenDevToolsWindow                     uintptr
		AddContainsFullScreenElementChanged    uintptr
		RemoveContainsFullScreenElementChanged uintptr
		GetContainsFullScreenElement           uintptr
		AddWebResourceRequested                uintptr
		RemoveWebResourceRequested             uintptr
		AddWebResourceRequestedFilter          uintptr
		RemoveWebResourceRequestedFilter       uintptr
		AddWindowCloseRequested                uintptr
		RemoveWindowCloseRequested             uintptr
	}
)

type (
	// ICoreWebView2Environment is the implementation of https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2environment
	ICoreWebView2Environment struct {
		Basic
		VTBL *ICoreWebView2EnvironmentVTBL
	}

	// ICoreWebView2EnvironmentVTBL is the implementation of https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2environment
	ICoreWebView2EnvironmentVTBL struct {
		BasicVTBL
		CreateCoreWebView2Controller     uintptr
		CreateWebResourceResponse        uintptr
		GetBrowserVersionString          uintptr
		AddNewBrowserVersionAvailable    uintptr
		RemoveNewBrowserVersionAvailable uintptr
	}
)

type (
	// ICoreWebView2Controller implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2controller?view=webview2-1.0.622.22
	ICoreWebView2Controller struct {
		VTBL *ICoreWebView2ControllerVTBL
	}

	// ICoreWebView2ControllerVTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2controller?view=webview2-1.0.622.22
	ICoreWebView2ControllerVTBL struct {
		BasicVTBL
		GetIsVisible                      uintptr
		PutIsVisible                      uintptr
		GetBounds                         uintptr
		PutBounds                         uintptr
		GetZoomFactor                     uintptr
		PutZoomFactor                     uintptr
		AddZoomFactorChanged              uintptr
		RemoveZoomFactorChanged           uintptr
		SetBoundsAndZoomFactor            uintptr
		MoveFocus                         uintptr
		AddMoveFocusRequested             uintptr
		RemoveMoveFocusRequested          uintptr
		AddGotFocus                       uintptr
		RemoveGotFocus                    uintptr
		AddLostFocus                      uintptr
		RemoveLostFocus                   uintptr
		AddAcceleratorKeyPressed          uintptr
		RemoveAcceleratorKeyPressed       uintptr
		GetParentWindow                   uintptr
		PutParentWindow                   uintptr
		NotifyParentWindowPositionChanged uintptr
		Close                             uintptr
		GetCoreWebView2                   uintptr
	}
)

type (
	// ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler is the implementation of https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2createcorewebview2environmentcompletedhandler.
	ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler struct {
		Basic
		VTBL *ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVTBL
	}

	// ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVTBL is the implementation of https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2createcorewebview2environmentcompletedhandler.
	ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVTBL struct {
		BasicVTBL
		// Invoke is ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerInvoke
		Invoke uintptr
	}

	// ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerInvoke: public HRESULT Invoke(HRESULT errorCode, ICoreWebView2Environment * createdEnvironment)
	ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerInvoke func(i *ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, p uintptr, createdEnvironment *ICoreWebView2Environment) uintptr
)

type (
	// ICoreWebView2CreateCoreWebView2ControllerCompletedHandler implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2createcorewebview2controllercompletedhandler
	ICoreWebView2CreateCoreWebView2ControllerCompletedHandler struct {
		Basic
		VTBL *ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerVTBL
	}

	// ICoreWebView2CreateCoreWebView2ControllerCompletedVTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2createcorewebview2controllercompletedhandler
	ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerVTBL struct {
		BasicVTBL
		Invoke uintptr
	}

	// ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerInvoke: public HRESULT Invoke(HRESULT errorCode, ICoreWebView2Controller * createdController)
	ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerInvoke func(i *ICoreWebView2CreateCoreWebView2ControllerCompletedHandler, p uintptr, createdController *ICoreWebView2Controller) uintptr
)