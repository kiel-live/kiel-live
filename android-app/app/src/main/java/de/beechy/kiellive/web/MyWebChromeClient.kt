package de.beechy.kiellive.web

import android.Manifest
import android.content.pm.PackageManager
import android.webkit.GeolocationPermissions
import android.webkit.WebChromeClient
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat
import de.beechy.kiellive.activities.MainActivity

class MyWebChromeClient(private val mainActivity: MainActivity) : WebChromeClient() {
    override fun onGeolocationPermissionsShowPrompt(
        origin: String?,
        callback: GeolocationPermissions.Callback
    ) {
        val perm = Manifest.permission.ACCESS_FINE_LOCATION
        if (ContextCompat.checkSelfPermission(
                mainActivity,
                perm
            ) == PackageManager.PERMISSION_GRANTED
        ) {
            callback.invoke(origin, true, false)
            return
        }

        if (!ActivityCompat.shouldShowRequestPermissionRationale(mainActivity, perm)) {
            // ask the user for permission
            ActivityCompat.requestPermissions(
                mainActivity,
                arrayOf(perm),
                MainActivity.REQUEST_FINE_LOCATION
            )

            // we will use these when user responds
            mainActivity.setGeolocationOrigin(origin)
            mainActivity.setGeolocationCallback(callback)
        }

    }
}
