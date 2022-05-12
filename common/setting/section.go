package setting

type SettingS struct {
	Live  LiveS
	Other OtherS
}

type LiveS struct {
	BiliRtmpUrl   string // 推流链接，由服务器地址和串流密钥拼接而成
	YoutubeUrl    string // Youtube 直播间或频道链接
	YoutubeCookie string // Youtube 登录 Cookie，用于限定直播等
}

type OtherS struct {
	Proxy      string // 代理地址
	FFmpegPath string // FFmpeg 路径
	FFmpegLogs bool   // 显示 ffmpeg 日志
	IdleInfo   bool   // 显示直播间等待信息
	CheckDelay int    // 直播间状态检查间隔
}

// 进行反序列化
func (s *Setting) ReadSection(v interface{}) error {
	err := s.vp.Unmarshal(v)
	if err != nil {
		return err
	}

	return nil
}
