package com.inkeliz.gowebview;

import android.os.Bundle;

import android.view.ViewGroup;
import android.app.Activity;
import android.view.View;
import android.view.KeyEvent;
import android.webkit.WebSettings;
import android.content.Context;
import android.os.PowerManager;
import android.webkit.WebViewClient;
import android.widget.Toast;
import android.webkit.WebView;
import android.util.Log;
import android.os.Build;
import android.os.Parcelable;
import android.net.Proxy;
import java.lang.reflect.*;
import android.util.ArrayMap;
import android.content.Intent;
import java.util.concurrent.Semaphore;
import android.net.http.SslError;
import android.webkit.SslErrorHandler;
import java.security.cert.Certificate;
import android.net.http.SslCertificate;
import java.security.PublicKey;
import java.security.cert.X509Certificate;
import java.security.MessageDigest;
import java.security.cert.CertificateFactory;
import java.io.ByteArrayInputStream;
import java.io.InputStream;
import android.os.Build.VERSION;
import android.os.Build.VERSION_CODES;
import android.util.Base64;
import android.os.Vibrator;
import android.os.VibrationEffect;

public class gowebview_android {
    private View primaryView;
    private WebView webBrowser;
    private PublicKey[] additionalCerts;

    public class gowebview_boolean {
        private boolean b;
        public void Set(Boolean r) {b = r;}
        public boolean Get() {return b;}
    }

    public class gowebview_webbrowser extends WebViewClient {
        @Override public void onReceivedSslError(WebView v, final SslErrorHandler sslHandler, SslError err){
            if (additionalCerts == null || additionalCerts.length == 0) {
                super.onReceivedSslError(v, sslHandler, err);
                return;
            }

            Certificate certificate = null;
            try{
                if (android.os.Build.VERSION.SDK_INT > android.os.Build.VERSION_CODES.Q) {
                      certificate = err.getCertificate().getX509Certificate();
                } else {
                    // Old APIs doesn't have such .getX509Certificate()
                    Bundle bundle = SslCertificate.saveState(err.getCertificate());
                    byte[] certificateBytes = bundle.getByteArray("x509-certificate");
                    if (certificateBytes != null) {
                        CertificateFactory certFactory = CertificateFactory.getInstance("X.509");
                        certificate = certFactory.generateCertificate(new ByteArrayInputStream(certificateBytes));
                    }
                }
            } catch (Exception e) {
                e.printStackTrace();
            }

            if (certificate == null) {
                super.onReceivedSslError(v, sslHandler, err);
                return;
            }

            for (int i = 0; i < additionalCerts.length; i++) {
                try{
                    certificate.verify(additionalCerts[i]);
                    sslHandler.proceed();
                    return;
                } catch (Exception e) {
                    e.printStackTrace();
                }
            }

            super.onReceivedSslError(v, sslHandler, err);
        }
    }

    // Executed when call `New(config *Config)`
    public void webview_create(View v) {
        primaryView = v;

        final Semaphore mutex = new Semaphore(0);

        ((Activity)primaryView.getContext()).runOnUiThread(new Runnable() {
            public void run() {
                webBrowser = new WebView(primaryView.getContext());
                WebSettings webSettings = webBrowser.getSettings();
                webSettings.setJavaScriptEnabled(true);
                if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.O_MR1) {
                    webSettings.setSafeBrowsingEnabled(false);
                }
                webSettings.setMixedContentMode(WebSettings.MIXED_CONTENT_COMPATIBILITY_MODE);
                webSettings.setUseWideViewPort(true);
                webSettings.setLoadWithOverviewMode(true);

                webBrowser.setWebViewClient(new gowebview_webbrowser());

                mutex.release();
            }
        });

        try {
            mutex.acquire();
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }

    // Executed when call `.SetURL(url string)`
    public void webview_navigate(String url) {
        ((Activity)primaryView.getContext()).runOnUiThread(new Runnable() {
            public void run() {
               webBrowser.loadUrl(url);
            }
        });
    }

    // Executed when call `.Run()` or `.SetVisibility()`
    public void webview_run() {
        ((Activity)primaryView.getContext()).runOnUiThread(new Runnable() {
            public void run() {
                ((Activity)primaryView.getContext()).setContentView(webBrowser);
            }
        });
    }

    public void webview_hibernate() {
        ((Activity)primaryView.getContext()).runOnUiThread(new Runnable() {
            public void run() {
                Intent intent= new Intent(Intent.ACTION_MAIN);
                intent.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
                intent.addCategory(Intent.CATEGORY_HOME);
                ((Activity)primaryView.getContext()).startActivity(intent);
            }
        });
    }

    public void webview_vibrate() {
        ((Activity)primaryView.getContext()).runOnUiThread(new Runnable() {
            public void run() {
                Vibrator vibrator = (Vibrator) ((Activity)primaryView.getContext()).getSystemService(Context.VIBRATOR_SERVICE);
    
                // Check whether device hardware supports vibration
                if (vibrator.hasVibrator()) {
                    // Vibrate for 500 milliseconds
                    if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                        // For newer versions use VibrationEffect
                        vibrator.vibrate(VibrationEffect.createOneShot(300, 255));  //VibrationEffect.DEFAULT_AMPLITUDE(5), the amplitude from 1 to 255
                    } else {
                        //deprecated in API 26
                        vibrator.vibrate(500);
                    }
                }
            }
        });
    }

    public void webview_wakelock() {
        PowerManager powerManager = (PowerManager) ((Activity)primaryView.getContext()).getSystemService(Context.POWER_SERVICE);
        powerManager.newWakeLock(PowerManager.PARTIAL_WAKE_LOCK, "goAndViewWakelock").acquire();
    }
    
    // Executed when call `.Destroy()`
    public void webview_destroy() {
        ((Activity)primaryView.getContext()).runOnUiThread(new Runnable() {
            public void run() {
                ((Activity)primaryView.getContext()).setContentView(primaryView);

                webBrowser.onPause();
                webBrowser.removeAllViews();
                webBrowser.pauseTimers();
                webBrowser.destroy();
            }
        });
    }

    // Executed when call `.SetVisibility()`
    public void webview_hide() {
        ((Activity)primaryView.getContext()).runOnUiThread(new Runnable() {
            public void run() {
                ((Activity)primaryView.getContext()).setContentView(primaryView);
            }
        });
    }

    public boolean webview_proxy(String host, String port) {
        final Semaphore mutex = new Semaphore(0);
        final gowebview_boolean result = new gowebview_boolean();

        ((Activity)primaryView.getContext()).runOnUiThread(new Runnable() {
            public void run() {
                Context app = webBrowser.getContext().getApplicationContext();

                System.setProperty("http.proxyHost", host);
                System.setProperty("http.proxyPort", port);
                System.setProperty("https.proxyHost", host);
                System.setProperty("https.proxyPort", port);

                try {
                    Field apk = app.getClass().getDeclaredField("mLoadedApk");
                    apk.setAccessible(true);

                    Field receivers = Class.forName("android.app.LoadedApk").getDeclaredField("mReceivers");
                    receivers.setAccessible(true);

                    for (Object map : ((ArrayMap) receivers.get(apk.get(app))).values()) {

                        for (Object receiver : ((ArrayMap) map).keySet()) {

                            Class<?> cls = receiver.getClass();
                            if (cls.getName().contains("ProxyChangeListener")) {

                                String proxyInfoName = "android.net.ProxyInfo";
                                if (Build.VERSION.SDK_INT <= Build.VERSION_CODES.KITKAT) {
                                    proxyInfoName = "android.net.ProxyProperties";
                                }

                                Intent intent = new Intent(Proxy.PROXY_CHANGE_ACTION);

                                Class proxyInfoClass = Class.forName(proxyInfoName);
                                if (proxyInfoClass != null) {
                                    Constructor proxyInfo = proxyInfoClass.getConstructor(String.class, Integer.TYPE, String.class);
                                    proxyInfo.setAccessible(true);
                                    intent.putExtra("proxy", (Parcelable) ((Object) proxyInfo.newInstance(host, Integer.parseInt(port), null)));
                                }

                                cls.getDeclaredMethod("onReceive", Context.class, Intent.class).invoke(receiver, app, intent);
                            }
                        }

                    }

                    result.Set(true);
                } catch(Exception e) {
                    e.printStackTrace();
                    result.Set(false);
                }

                mutex.release();
            }
        });

        try {
            mutex.acquire();
        } catch (InterruptedException e) {
            e.printStackTrace();
        }

        return result.Get();
    }

    public boolean webview_certs(String der) {
        String[] sCerts = der.split(";");

        additionalCerts = new PublicKey[sCerts.length];

        for (int i = 0; i < sCerts.length; i++) {
            InputStream streamCert = new ByteArrayInputStream(Base64.decode(sCerts[i], android.util.Base64.DEFAULT));

            try {
                CertificateFactory factory = CertificateFactory.getInstance("X.509");
                 X509Certificate cert = (X509Certificate)factory.generateCertificate(streamCert);

                 additionalCerts[i] = cert.getPublicKey();
            } catch(Exception e) {
                e.printStackTrace();
                return false;
            }
        }

        return true;
    }
}
