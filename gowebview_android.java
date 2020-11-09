package com.inkeliz.gowebview;

import android.os.Bundle;

import android.view.ViewGroup;
import android.app.Activity;
import android.view.View;
import android.view.KeyEvent;
import android.webkit.WebSettings;
import android.content.Context;
import android.webkit.WebViewClient;
import android.widget.Toast;
import android.webkit.WebView;
import android.util.Log;

public class gowebview_android {
    private View primaryView;
    private WebView webBrowser;

    // Executed when call `New(config *Config)`
    public void webview_create(View v) {
        primaryView = v;

        ((Activity)primaryView.getContext()).runOnUiThread(new Runnable() {
            public void run() {
                webBrowser = new WebView(primaryView.getContext());
                WebSettings webSettings = webBrowser.getSettings();
                webSettings.setJavaScriptEnabled(true);
                webSettings.setSafeBrowsingEnabled(false);
                webSettings.setMixedContentMode(WebSettings.MIXED_CONTENT_COMPATIBILITY_MODE);
                webSettings.setUseWideViewPort(true);
                webSettings.setLoadWithOverviewMode(true);

                webBrowser.setWebViewClient(new WebViewClient());
            }
        });
    }

    // Executed when call `.SetURL(url string)`
    public void webview_navigate(String url) {
        webBrowser.loadUrl(url);
    }

    // Executed when call `.Run()`
    public void webview_run() {
        ((Activity)primaryView.getContext()).runOnUiThread(new Runnable() {
            public void run() {
                ((Activity)primaryView.getContext()).setContentView(webBrowser);
            }
        });
    }

    // Executed when call `.Destroy()`
    public void webview_destroy() {
        ((Activity)primaryView.getContext()).runOnUiThread(new Runnable() {
            public void run() {
                webBrowser.clearHistory();
                webBrowser.clearCache(true);
                webBrowser.onPause();
                webBrowser.removeAllViews();
                webBrowser.pauseTimers();
                webBrowser.destroy();

                ((Activity)primaryView.getContext()).setContentView(primaryView);
            }
        });
    }
}