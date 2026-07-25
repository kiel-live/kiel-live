package de.beechy.kiellive.activities

import android.content.Intent
import android.os.Bundle
import android.os.Handler
import android.view.View
import android.widget.RelativeLayout
import androidx.appcompat.app.AppCompatActivity
import de.beechy.kiellive.R

class SplashActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.splash_activity)

        // wait some time before opening web-view
        val handler = Handler()
        handler.postDelayed(Runnable { this.openMain() }, 500)

        // add click-listener to skip splash-screen
        val layout = findViewById<RelativeLayout>(R.id.layout)
        layout.setOnClickListener(View.OnClickListener { v: View? ->
            // stop handler
            handler.removeCallbacksAndMessages(null)
            openMain()
        })
    }

    private fun openMain() {
        val intent = Intent(this@SplashActivity, MainActivity::class.java)
        startActivity(intent)
        finish()
    }
}