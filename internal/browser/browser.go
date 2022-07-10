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
		Set("--autoplay-policy", "user-gesture-required").
		Set("--disable-component-update").
		Set("--disable-domain-reliability").
		Set("--disable-features", "AudioServiceOutOfProcess,IsolateOrigins,site-per-process").
		Set("--disable-print-preview").
		Set("--disable-setuid-sandbox").
		Set("--disable-speech-api").
		Set("--disable-web-security").
		Set("--disk-cache-size", "33554432").
		Set("--enable-features", "SharedArrayBuffer").
		Set("--hide-scrollbars").
		Set("--ignore-gpu-blocklist").
		Set("--in-process-gpu").
		Set("--mute-audio").
		Set("--no-default-browser-check").
		Set("--no-pings").
		Set("--no-sandbox").
		Set("--no-zygote").
		Set("--use-gl", "angle").
		Set("--use-angle", "swiftshader").
		Set("--single-process").
		MustLaunch()

	return rod.New().ControlURL(u).MustConnect()
}
