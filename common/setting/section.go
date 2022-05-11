package setting

type LiveSettingS struct {
	BiliRtmpUrl   string // 推流链接，由服务器地址和串流密钥拼接而成
	YoutubeUrl    string // Youtube 直播间或频道链接
	YoutubeCookie string // Youtube 登录 Cookie，用于限定直播等
	Proxy         string // 代理地址
}

// 进行反序列化
func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	return nil
}
