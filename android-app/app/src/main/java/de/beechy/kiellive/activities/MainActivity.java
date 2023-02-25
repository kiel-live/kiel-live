package de.beechy.kiellive.activities;

import androidx.annotation.NonNull;
import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.content.pm.PackageManager;
import android.net.Uri;
import android.os.Build;
import android.os.Bundle;
import android.view.View;
import android.webkit.GeolocationPermissions;
import android.webkit.WebSettings;
import android.webkit.WebView;

import java.util.HashMap;
import java.util.Map;

import de.beechy.kiellive.Config;
import de.beechy.kiellive.web.MyWebChromeClient;
import de.beechy.kiellive.web.MyWebViewClient;

public class MainActivity extends AppCompatActivity {

    public static final int REQUEST_FINE_LOCATION = 1;
    private WebView myWebView;
    private String geolocationOrigin;
    private GeolocationPermissions.Callback geolocationCallback;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        // set special start path if activity got url with intent
        String path = getStartPath(getIntent());
        if (path == null) {
            path = "/";
        }

        myWebView = new WebView(this.getBaseContext());
        setContentView(myWebView);

        WebSettings webSettings = myWebView.getSettings();
        webSettings.setJavaScriptEnabled(true);
        webSettings.setCacheMode(WebSettings.LOAD_DEFAULT);
        webSettings.setDatabaseEnabled(true);
        webSettings.setDomStorageEnabled(true);
        webSettings.setGeolocationEnabled(true);
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            webSettings.setSafeBrowsingEnabled(false);
        }

        MyWebViewClient webViewClient = new MyWebViewClient(this);

        myWebView.setLayerType(View.LAYER_TYPE_HARDWARE, null);
        myWebView.setOnApplyWindowInsetsListener((view, insets) -> {
            webViewClient.setStatusBarHeight(insets.getStableInsetTop());
            return insets.consumeSystemWindowInsets();
        });

        myWebView.setWebViewClient(webViewClient);
        myWebView.setWebChromeClient(new MyWebChromeClient(this));

        Map<String, String> extraHeaders = new HashMap<>();
        extraHeaders.put("Referer", "android-app://" + getApplication().getPackageName());
        myWebView.loadUrl(Config.APP_URL + path, extraHeaders);
    }

    @Override
    public void onBackPressed(){
        if (myWebView.canGoBack()) {
            myWebView.goBack();
        } else {
            super.onBackPressed();
        }
    }

    @Override
    public void onRequestPermissionsResult(int requestCode, @NonNull String[] permissions, @NonNull int[] grantResults) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults);
        switch (requestCode) {
            case REQUEST_FINE_LOCATION:
                boolean allow = false;
                if (grantResults[0] == PackageManager.PERMISSION_GRANTED) {
                    // user has allowed this permission
                    allow = true;
                }
                if (geolocationCallback != null) {
                    // call back to web chrome client
                    geolocationCallback.invoke(geolocationOrigin, allow, false);
                }
                break;
        }
    }

    public void setGeolocationOrigin(String geolocationOrigin) {
        this.geolocationOrigin = geolocationOrigin;
    }

    public void setGeolocationCallback(GeolocationPermissions.Callback geolocationCallback) {
        this.geolocationCallback = geolocationCallback;
    }

    private String getStartPath(Intent intent) {
        if (intent == null) {
            return null;
        }

        Uri uri = intent.getData();

        if (uri == null) {
            return null;
        }

        if (!uri.getHost().equals(Config.APP_HOST)) {
            return null;
        }

        if (uri.getPath() == null) {
            return null;
        }

        return uri.getPath();
    }
}
