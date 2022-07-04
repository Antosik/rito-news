package contentstack

import (
	"fmt"
	"net/http"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type ContentStackKeys struct {
	api_key      string
	access_token string
}

func (keys ContentStackKeys) String() string {
	return fmt.Sprintf("api_key: %v, access_token: %v", keys.api_key, keys.access_token)
}

func GetContentStackKeys(url string, selectorToWait string, params *ContentStackQueryParameters) *ContentStackKeys {
	var keys ContentStackKeys

	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	router := browser.HijackRequests()
	defer router.MustStop()

	target := fmt.Sprintf("https://cdn.contentstack.io/v3/content_types/%s/*", params.ContentType)

	router.MustAdd(target, func(ctx *rod.Hijack) {
		headers := ctx.Request.Headers()

		at, at_ok := headers["access_token"]
		if at_ok {
			(&keys).access_token = at.Str()
		}

		ak, ak_ok := headers["api_key"]
		if ak_ok {
			(&keys).api_key = ak.Str()
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
