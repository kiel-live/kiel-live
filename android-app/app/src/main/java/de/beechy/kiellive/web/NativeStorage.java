package de.beechy.kiellive.web;

import android.content.Context;
import android.content.SharedPreferences;
import android.webkit.JavascriptInterface;

/**
 * localStorage-like key/value store backed by SharedPreferences.
 * Exposed to the WebView as `window.NativeStorage` so web data survives
 * a change of the WebView origin / domain.
 */
public class NativeStorage {
    public static final String JS_NAME = "NativeStorage";
    private static final String PREFS_NAME = "web_storage";

    private final SharedPreferences prefs;

    public NativeStorage(Context context) {
        prefs = context.getApplicationContext().getSharedPreferences(PREFS_NAME, Context.MODE_PRIVATE);
    }

    @JavascriptInterface
    public String getItem(String key) {
        return prefs.getString(key, null);
    }

    @JavascriptInterface
    public void setItem(String key, String value) {
        prefs.edit().putString(key, value).apply();
    }

    @JavascriptInterface
    public void removeItem(String key) {
        prefs.edit().remove(key).apply();
    }
}
