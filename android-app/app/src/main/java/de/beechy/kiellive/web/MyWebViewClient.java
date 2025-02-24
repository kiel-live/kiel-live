package de.beechy.kiellive.web;

import android.content.Context;
import android.content.Intent;
import android.net.Uri;
import android.webkit.WebView;
import android.webkit.WebViewClient;

import de.beechy.kiellive.Config;

public class MyWebViewClient extends WebViewClient {
    private Context context;
    private int statusBarHeight;

    public MyWebViewClient(Context context) {
        this.context = context;
    }

    @Override
    public boolean shouldOverrideUrlLoading(WebView view, String url) {
        if (Config.APP_HOST.equals(Uri.parse(url).getHost())) {
            // This is my website, so do not override; let my WebView load the page
            return false;
        }

        // Otherwise, the link is not for a page on my site, so launch another Activity that handles URLs
        Intent intent = new Intent(Intent.ACTION_VIEW, Uri.parse(url));
        context.startActivity(intent);
        return true;
    }

    // https://stackoverflow.com/a/30270803
    private final static String CREATE_CUSTOM_SHEET =
            "if (typeof(document.head) != 'undefined' && typeof(customSheet) == 'undefined') {"
                    + "var customSheet = (function() {"
                    + "var style = document.createElement(\"style\");"
                    + "style.appendChild(document.createTextNode(\"\"));"
                    + "document.head.appendChild(style);"
                    + "return style.sheet;"
                    + "})();"
                    + "}";

    // https://stackoverflow.com/a/30270803
    void injectCssIntoWebView(WebView webView, int statusBarHeight) {
        StringBuilder jsUrl = new StringBuilder("javascript:");
        jsUrl.append(CREATE_CUSTOM_SHEET)
            .append("if (typeof(customSheet) != 'undefined') {")
            .append("const pixel = " + statusBarHeight + " / window.devicePixelRatio + 8;")
            .append("customSheet.insertRule('#app-bar,#settings-container,#details-popup { margin-top: ' + pixel + 'px; }', 0);")
            .append("}");

        webView.loadUrl(jsUrl.toString());
    }

    public void setStatusBarHeight(int statusBarHeight) {
        this.statusBarHeight = statusBarHeight;
    }

    @Override
    public void onPageFinished(WebView webView, String url) {
        injectCssIntoWebView(
                webView,
                statusBarHeight
        );
    }
}
