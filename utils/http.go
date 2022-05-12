package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func HttpGet(url string, headers map[string]string) ([]byte, error) {
	client := &http.Client{}

	reqest, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	if headers != nil {
		for key, val := range headers {
			reqest.Header.Set(key, val)
		}
	}

	response, err := client.Do(reqest)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP status code: %d", response.StatusCode)
	}

	respByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return respByte, nil
}

/*
// 将rawurl 转换为 url.URL
func GetProxyUrl(proxy string) *url.URL {
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		return nil
	}
	return proxyUrl
}
*/
