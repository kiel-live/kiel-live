package de.beechy.kiellive.activities;

import androidx.annotation.Nullable;
import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.os.Handler;
import android.view.View;
import android.widget.RelativeLayout;

import de.beechy.kiellive.R;

public class SplashActivity extends AppCompatActivity {

    @Override
    protected void onCreate(@Nullable Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.splash_activity);

        // wait some time before opening web-view
        Handler handler = new Handler();
        handler.postDelayed(this::openMain, 500);

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