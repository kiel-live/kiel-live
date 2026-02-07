package de.beechy.kiellive.web;

import android.content.Context;
import android.content.Intent;
import android.net.Uri;
import android.webkit.WebView;
import android.webkit.WebViewClient;

import de.beechy.kiellive.Config;

public class MyWebViewClient extends WebViewClient {
    private final Context context;
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

    /**
     * Injects the safe area top inset as a CSS custom property (CSS variable).
     * https://stackoverflow.com/a/30270803
     */
    void injectCssIntoWebView(WebView webView, int statusBarHeight) {
        // Calculate the pixel value accounting for device pixel ratio
        String jsCode = "javascript:{"
                + "const pixel = " + statusBarHeight + " / window.devicePixelRatio;"
                + "document.documentElement.style.setProperty('--safe-area-top', pixel + 'px');"
                + "}";

        webView.loadUrl(jsCode);
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
