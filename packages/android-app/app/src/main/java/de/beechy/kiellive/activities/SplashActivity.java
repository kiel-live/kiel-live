package de.beechy.kiellive.activities;

import androidx.annotation.Nullable;
import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.content.pm.PackageInfo;
import android.content.pm.PackageManager;
import android.os.Build;
import android.os.Bundle;
import android.os.Handler;
import android.transition.Explode;
import android.view.View;
import android.view.Window;
import android.widget.RelativeLayout;
import android.widget.TextView;

import de.beechy.kiellive.R;

public class SplashActivity extends AppCompatActivity {

    @Override
    protected void onCreate(@Nullable Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.splash_activity);

        // set current version
        try {
            PackageInfo pInfo = getPackageManager().getPackageInfo(getPackageName(), 0);
            String appVersion = pInfo.versionName;
            TextView appVersionText = findViewById(R.id.app_version);
            appVersionText.setText("Version " + appVersion);
        } catch (PackageManager.NameNotFoundException e) {
            e.printStackTrace();
        }

        // wait some time before opening web-view
        Handler handler = new Handler();
        handler.postDelayed(() -> {
            openMain();
        },1000);

        // add click-listener to skip splash-screen
        RelativeLayout layout = findViewById(R.id.layout);
        layout.setOnClickListener((View v) -> {
            // stop handler
            handler.removeCallbacksAndMessages(null);
            openMain();
        });
    }
    private void openMain() {
        Intent intent = new Intent(SplashActivity.this, MainActivity.class);
        startActivity(intent);
        finish();
    }
}