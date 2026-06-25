package de.beechy.kiellive.activities

import android.annotation.SuppressLint
import android.content.Intent
import android.content.pm.PackageManager
import android.os.Build
import android.os.Bundle
import android.view.View
import android.webkit.GeolocationPermissions
import android.webkit.WebSettings
import android.webkit.WebView
import androidx.activity.OnBackPressedCallback
import androidx.appcompat.app.AppCompatActivity
import de.beechy.kiellive.BuildConfig
import de.beechy.kiellive.Config
import de.beechy.kiellive.web.MyWebChromeClient
import de.beechy.kiellive.web.MyWebViewClient

class MainActivity : AppCompatActivity() {
    private var myWebView: WebView? = null
    private var geolocationOrigin: String? = null
    private var geolocationCallback: GeolocationPermissions.Callback? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        // set special start path if activity got url with intent
        var path = getStartPath(intent)
        if (path == null) {
            path = "/"
        }

        myWebView = WebView(this)
        setContentView(myWebView)

        val webSettings = myWebView!!.settings
        webSettings.javaScriptEnabled = true
        webSettings.cacheMode = WebSettings.LOAD_DEFAULT
        webSettings.domStorageEnabled = true
        webSettings.setGeolocationEnabled(true)
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            webSettings.safeBrowsingEnabled = false
        }

        val webViewClient = MyWebViewClient(this)

        myWebView!!.setLayerType(View.LAYER_TYPE_HARDWARE, null)
        myWebView!!.webViewClient = webViewClient
        myWebView!!.webChromeClient = MyWebChromeClient(this)

        val extraHeaders: MutableMap<String?, String?> = HashMap()
        extraHeaders["Referer"] = "android-app://" + application.packageName
        extraHeaders["app-name"] = application.packageName
        extraHeaders["app-version"] = BuildConfig.VERSION_NAME
        myWebView!!.loadUrl(Config.APP_URL + path, extraHeaders)

        val callback: OnBackPressedCallback = object : OnBackPressedCallback(true) {
            override fun handleOnBackPressed() {
                if (myWebView!!.canGoBack()) {
                    myWebView!!.goBack()
                } else {
                    isEnabled = false
                    this@MainActivity.onBackPressedDispatcher.onBackPressed()
                }
            }
        }
        onBackPressedDispatcher.addCallback(this, callback)
    }

    override fun onRequestPermissionsResult(
        requestCode: Int,
        permissions: Array<String?>,
        grantResults: IntArray
    ) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults)
        when (requestCode) {
            REQUEST_FINE_LOCATION -> {
                var allow = false
                if (grantResults[0] == PackageManager.PERMISSION_GRANTED) {
                    // user has allowed this permission
                    allow = true
                }
                if (geolocationCallback != null) {
                    // call back to web chrome client
                    geolocationCallback!!.invoke(geolocationOrigin, allow, false)
                }
            }
        }
    }

    fun setGeolocationOrigin(geolocationOrigin: String?) {
        this.geolocationOrigin = geolocationOrigin
    }

    fun setGeolocationCallback(geolocationCallback: GeolocationPermissions.Callback?) {
        this.geolocationCallback = geolocationCallback
    }

    private fun getStartPath(intent: Intent?): String? {
        if (intent == null) {
            return null
        }

        val uri = intent.data ?: return null

        if (uri.host != Config.APP_HOST) {
            return null
        }

        if (uri.path == null) {
            return null
        }

        return uri.path
    }

    companion object {
        const val REQUEST_FINE_LOCATION: Int = 1
    }
}
