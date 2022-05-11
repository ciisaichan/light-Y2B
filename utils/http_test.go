package utils

import (
	"testing"

	_ "github.com/ciisaichan/light-Y2B/bootstrap"
	"github.com/ciisaichan/light-Y2B/global"
)

func TestHttp(t *testing.T) {
	t.Log(global.LiveSetting.Proxy)
	body, err := HttpGet("http://www.google.com", nil)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(body))
}
