package youtube

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ciisaichan/light-Y2B/utils"
)

var hlsvpRegexp *regexp.Regexp
var httpUA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36 Edg/101.0.1210.39"

func init() {
	hlsvpRegexp = regexp.MustCompile(`(\\\\\\\"hlsManifestUrl\\\\\\\":\\\\\\\"(.+?)\\\\\\\")|(\"hlsManifestUrl\":\"(.+?)\")`)
}

func GetLiveStreamURL(url string, cookie string) (string, error) {
	var heads = make(map[string]string)
	heads["User-Agent"] = httpUA
	if cookie != "" {
		heads["Cookie"] = cookie
	}

	liveHtmlByte, err := utils.HttpGet(url, heads)
	if err != nil {
		return "", fmt.Errorf("get live page failed: %s", err)
	}
	liveHtml := string(liveHtmlByte)

	params := hlsvpRegexp.FindStringSubmatch(liveHtml)
	if len(params) == 0 {
		return "", fmt.Errorf("hlsvp url not found")
	}

	if !strings.Contains(liveHtml, `"isLive":true`) {
		return "", fmt.Errorf("not live")
	}

	var hlsvpMatchUrl string
	if params[2] == "" {
		hlsvpMatchUrl = params[4]
	} else {
		hlsvpMatchUrl = params[2]
	}

	m3u8ListByte, err := utils.HttpGet(hlsvpMatchUrl, heads)
	if err != nil {
		return "", fmt.Errorf("get m3u8 list failed: %s", err)
	}

	m3u8List := strings.Split(string(m3u8ListByte), "\n")
	var m3u8Url string
	for i := len(m3u8List) - 1; i >= 0; i-- {
		if strings.HasPrefix(m3u8List[i], "#EXT-X-STREAM-INF") {
			m3u8Url = m3u8List[i+1]
			break
		}
	}

	if m3u8Url != "" {
		return m3u8Url, nil
	} else {
		return "", fmt.Errorf("m3u8 url not found")
	}
}
