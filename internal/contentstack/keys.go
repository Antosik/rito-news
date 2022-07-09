package contentstack

import (
	"fmt"
	"net/http"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type Keys struct {
	apiKey      string
	accessToken string
}

func (keys Keys) String() string {
	return fmt.Sprintf("api_key: %v, access_token: %v", keys.apiKey, keys.accessToken)
}

func GetKeys(url string, selectorToWait string, params *Parameters) *Keys {
	var keys Keys

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
		Headless(true).
		MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	router := browser.HijackRequests()
	defer router.MustStop()

	target := fmt.Sprintf("https://cdn.contentstack.io/v3/content_types/%s/*", params.ContentType)

	router.MustAdd(target, func(ctx *rod.Hijack) {
		headers := ctx.Request.Headers()

		at, atOk := headers["access_token"]
		if atOk {
			(&keys).accessToken = at.Str()
		}

		ak, akOk := headers["api_key"]
		if akOk {
			(&keys).apiKey = ak.Str()
		}

		_ = ctx.LoadResponse(http.DefaultClient, true)
	})

	go router.Run()

	page := browser.MustPage(url)
	defer page.Close()

	if selectorToWait == "body" {
		page.MustWaitLoad()
	} else {
		page.MustElement(selectorToWait)
	}

	// Hijack requests under the scope of a page
	page.HijackRequests()

	return &keys
}
