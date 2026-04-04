package de.beechy.kiellive.web

import android.content.Context
import android.content.Intent
import android.webkit.WebView
import android.webkit.WebViewClient
import androidx.core.net.toUri
import de.beechy.kiellive.Config

class MyWebViewClient(private val context: Context) : WebViewClient() {
    @Deprecated("Deprecated in Java")
    override fun shouldOverrideUrlLoading(view: WebView?, url: String?): Boolean {
        if (Config.APP_HOST == url?.toUri()?.host) {
            // This is my website, so do not override; let my WebView load the page
            return false
        }

        // Otherwise, the link is not for a page on my site, so launch another Activity that handles URLs
        val intent = Intent(Intent.ACTION_VIEW, url?.toUri())
        context.startActivity(intent)
        return true
    }
}
