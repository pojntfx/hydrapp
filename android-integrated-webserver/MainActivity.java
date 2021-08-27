package com.pojtinger.felicitas.integratedWebserverExample;

import android.app.Activity;
import android.content.Intent;
import android.net.Uri;
import android.os.Bundle;
import android.view.KeyEvent;
import android.webkit.PermissionRequest;
import android.webkit.ValueCallback;
import android.webkit.WebChromeClient;
import android.webkit.WebResourceRequest;
import android.webkit.WebSettings;
import android.webkit.WebView;
import android.webkit.WebViewClient;
import android.widget.LinearLayout;

public class MainActivity extends Activity {
    static {
        System.loadLibrary("backend");
    }

    private native String LaunchBackend();

    private ValueCallback<Uri[]> fileChooserCallback;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        Uri home = Uri.parse(LaunchBackend());

        WebView view = new WebView(this);
        view.setLayoutParams(new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT,
                LinearLayout.LayoutParams.MATCH_PARENT));

        WebSettings settings = view.getSettings();
        settings.setAllowContentAccess(true);
        settings.setAllowFileAccess(true);
        settings.setDatabaseEnabled(true);
        settings.setDomStorageEnabled(true);
        settings.setGeolocationEnabled(true);
        settings.setJavaScriptCanOpenWindowsAutomatically(true);
        settings.setJavaScriptEnabled(true);
        settings.setLoadsImagesAutomatically(true);
        settings.setMediaPlaybackRequiresUserGesture(false);
        settings.setSupportMultipleWindows(true);

        view.setWebViewClient(new WebViewClient() {
            @Override
            public boolean shouldOverrideUrlLoading(WebView vw, WebResourceRequest request) {
                if (request.getUrl().toString().contains(home.getHost())) {
                    vw.loadUrl(request.getUrl().toString());
                } else {
                    Intent intent = new Intent(Intent.ACTION_VIEW, request.getUrl());
                    vw.getContext().startActivity(intent);
                }

                return true;
            }
        });
        view.setWebChromeClient(new WebChromeClient() {
            @Override
            public void onPermissionRequest(final PermissionRequest request) {
                request.grant(request.getResources());
            }

            @Override
            public boolean onShowFileChooser(WebView vw, ValueCallback<Uri[]> filePathCallback,
                    FileChooserParams fileChooserParams) {
                if (fileChooserCallback != null) {
                    fileChooserCallback.onReceiveValue(null);
                }
                fileChooserCallback = filePathCallback;

                Intent selectionIntent = new Intent(Intent.ACTION_GET_CONTENT);
                selectionIntent.addCategory(Intent.CATEGORY_OPENABLE);
                selectionIntent.setType("*/*");

                Intent chooserIntent = new Intent(Intent.ACTION_CHOOSER);
                chooserIntent.putExtra(Intent.EXTRA_INTENT, selectionIntent);
                startActivityForResult(chooserIntent, 0);

                return true;
            }
        });
        view.setOnKeyListener((v, keyCode, event) -> {
            WebView vw = (WebView) v;
            if (event.getAction() == KeyEvent.ACTION_DOWN && keyCode == KeyEvent.KEYCODE_BACK && vw.canGoBack()) {
                vw.goBack();

                return true;
            }

            return false;
        });
        view.setDownloadListener((uri, userAgent, contentDisposition, mimetype, contentLength) -> handleURI(uri));
        view.setOnLongClickListener(v -> {
            handleURI(((WebView) v).getHitTestResult().getExtra());

            return true;
        });

        view.loadUrl(home.toString());

        setContentView(view);
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent intent) {
        super.onActivityResult(requestCode, resultCode, intent);

        fileChooserCallback.onReceiveValue(new Uri[] { Uri.parse(intent.getDataString()) });
        fileChooserCallback = null;
    }

    private void handleURI(String uri) {
        if (uri != null) {
            Intent i = new Intent(Intent.ACTION_VIEW);
            i.setData(Uri.parse(uri.replaceFirst("^blob:", "")));

            startActivity(i);
        }
    }
}