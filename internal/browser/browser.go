// Package browser is the utility package that allows to launch browser with the perfomance oriented arguments
package browser

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func NewBrowser() *rod.Browser {
	path, _ := launcher.LookPath()
	u := launcher.New().
		Bin(path).
		Set("--allow-running-insecure-content").
		Set("--autoplay-policy", "user-gesture-required"). // Don't render video
		Set("--disable-component-update").
		Set("--disable-domain-reliability"). // Disables Domain Reliability Monitoring.
		Set("--disable-features", "AudioServiceOutOfProcess,IsolateOrigins,site-per-process").
		Set("--disable-extensions").     // Disable all chrome extensions
		Set("--disable-notifications").  // Disables the Web Notification and the Push APIs.
		Set("--disable-setuid-sandbox"). // Disable the setuid sandbox (Linux only).
		Set("--disable-speech-api").     // Disables the Web Speech API (both speech recognition and synthesis).
		Set("--disable-web-security").   // Don't enforce the same-origin policy.
		Set("--disk-cache-size", "33554432").
		Set("--disable-gpu").
		Set("--blink-settings", "imagesEnabled=false").
		Set("--enable-automation"). // Disable a few things considered not appropriate for automation.
		Set("--enable-features", "SharedArrayBuffer").
		Set("--headless").                 // Run in headless mode
		Set("--hide-scrollbars").          // Hide scrollbars
		Set("--ignore-gpu-blocklist").     // Ignores GPU blocklist.
		Set("--in-process-gpu").           // Run the GPU process as a thread in the browser process.
		Set("--mute-audio").               // Mute any audio
		Set("--no-default-browser-check"). // Disables the default browser check.
		Set("--no-first-run").             // Skip first run wizards
		Set("--no-pings").                 // Don't send hyperlink auditing pings
		Set("--no-sandbox").               // Disables the sandbox for all process types that are normally sandboxed.
		Set("--no-zygote").                // Disables the use of a zygote process for forking child processes.
		Set("--use-gl", "angle").          // Select which implementation of GL the GPU process should use.
		Set("--use-angle", "swiftshader"). // Select which ANGLE backend to use.
		Set("--single-process").           // Runs the renderer and plugins in the same process as the browser.
		MustLaunch()

	return rod.New().ControlURL(u).MustConnect()
}
