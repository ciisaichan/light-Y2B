package youtube

import (
	"strings"

	"github.com/ciisaichan/light-Y2B/utils"
)

func IsLiving(url string, cookie string) (bool, error) {
	var heads = make(map[string]string)
	heads["User-Agent"] = httpUA
	if cookie != "" {
		heads["Cookie"] = cookie
	}

	liveHtmlByte, err := utils.HttpGet(url, heads)
	if err != nil {
		return false, err
	}

	if strings.Contains(string(liveHtmlByte), `"isLive":true`) {
		return true, nil
	}

	return false, nil
}

func IsChannelURL(url string) bool {
	return strings.Contains(url, "/channel/")
}
