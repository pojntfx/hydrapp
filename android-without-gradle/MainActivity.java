package com.example.helloandroid;

import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.net.Uri;
import android.os.Bundle;
import android.webkit.WebView;
import android.widget.FrameLayout;
import android.widget.LinearLayout;

public class MainActivity extends Activity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        Context context = getApplicationContext();
        super.onCreate(savedInstanceState);

        LinearLayout layout = new LinearLayout(this);
        layout.setLayoutParams(new FrameLayout.LayoutParams(FrameLayout.LayoutParams.MATCH_PARENT, FrameLayout.LayoutParams.MATCH_PARENT));

        WebView webView = new WebView(this);
        webView.setLayoutParams(new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.MATCH_PARENT));
        webView.getSettings().setJavaScriptEnabled(true);

        Intent intent = getIntent();
        Uri data = intent.getData();

        // webView.loadUrl(data.toString());
        webView.loadUrl("https://go-app.dev/");
        layout.addView(webView);

        setContentView(layout);
    }
}