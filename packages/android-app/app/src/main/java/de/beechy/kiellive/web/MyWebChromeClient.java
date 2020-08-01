package de.beechy.kiellive.web;

import android.Manifest;
import android.app.Activity;
import android.content.pm.PackageManager;
import android.os.Build;
import android.webkit.GeolocationPermissions;
import android.webkit.WebChromeClient;

import androidx.core.app.ActivityCompat;
import androidx.core.content.ContextCompat;

import de.beechy.kiellive.activities.MainActivity;

public class MyWebChromeClient extends WebChromeClient {
    private MainActivity mainActivity;

    public MyWebChromeClient(MainActivity mainActivity) {
        this.mainActivity = mainActivity;
    }

    @Override
    public void onGeolocationPermissionsShowPrompt(String origin, GeolocationPermissions.Callback callback) {
        // Geolocation permissions coming from this app's Manifest will only be valid for devices with
        // API_VERSION < 23. On API 23 and above, we must check for permissions, and possibly
        // ask for them.
        String perm = Manifest.permission.ACCESS_FINE_LOCATION;
        if (Build.VERSION.SDK_INT < Build.VERSION_CODES.M ||
                ContextCompat.checkSelfPermission(mainActivity, perm) == PackageManager.PERMISSION_GRANTED) {
            // we're on SDK < 23 OR user has already granted permission
            callback.invoke(origin, true, false);
        } else {
            if (!ActivityCompat.shouldShowRequestPermissionRationale(mainActivity, perm)) {
                // ask the user for permission
                ActivityCompat.requestPermissions(mainActivity, new String[] {perm}, mainActivity.REQUEST_FINE_LOCATION);

                // we will use these when user responds
                mainActivity.setGeolocationOrigin(origin);
                mainActivity.setGeolocationCallback(callback);
            }
        }
    }
}
