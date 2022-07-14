package contentstack

import (
	"fmt"
	"net/http"

	"github.com/Antosik/rito-news/internal/browser"
	"github.com/go-rod/rod"
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

	browser := browser.NewBrowser()
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

	return &keys
}
