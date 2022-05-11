package main

import (
	"context"
	"flag"
	"fmt"
	"light-Y2B/ffmpeg"
	"light-Y2B/global"
	"light-Y2B/logger"
	"light-Y2B/platform/youtube"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	ytLiveUrl     string
	ytCookie      string
	pushURL       string
	ffmpegPath    string
	checkDelay    int
	hideIdleInfo  bool
	hideFFmpegLog bool

	osChannel  chan os.Signal
	ctx        context.Context
	mainCancel context.CancelFunc
	mainWg     sync.WaitGroup
)

func init() {
	flag.StringVar(&ytLiveUrl, "yt-url", "", "Youtube 直播间或频道链接")
	flag.StringVar(&ytCookie, "yt-cookie", "", "Youtube 登录 Cookie，用于限定直播等")
	flag.StringVar(&pushURL, "push-url", "", "推流链接，由服务器地址和串流密钥拼接而成")
	flag.StringVar(&ffmpegPath, "ff-path", "ffmpeg", "ffmpeg 路径，默认从环境变量获取")
	flag.IntVar(&checkDelay, "ck-delay", 10, "直播间检查间隔，单位：秒，默认 10 秒")
	flag.BoolVar(&hideIdleInfo, "hide-idle", false, "是否隐藏等待信息，默认 不隐藏")
	flag.BoolVar(&hideFFmpegLog, "hide-fflog", false, "是否隐藏 ffmpeg 日志信息，默认 不隐藏")
}

func main() {
	fmt.Printf("# Light-Y2B [%s]\n", global.Version)
	fmt.Println("# 基于 FFmpeg 的轻量级 Youtube 转播程序\n")
	flag.Parse()

	if ytLiveUrl == "" || pushURL == "" {
		flag.Usage()
		return
	}

	//os.Setenv("http_proxy", "http://127.0.0.1:8080")
	//os.Setenv("https_proxy", "http://127.0.0.1:8080")

	osChannel = make(chan os.Signal, 1)
	signal.Notify(osChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	ctx, mainCancel = context.WithCancel(context.Background())

	mainWg.Add(1)

	go func() {
		for {
			select {
			case <-osChannel:
				global.Aborting = true
				mainCancel()
				mainWg.Done()
			}
		}
	}()

	if youtube.IsChannelURL(ytLiveUrl) {
		ytLiveUrl = ytLiveUrl + "/live"
	}

	go handleLoop()
	logger.L.Info("开始监听直播间状态...")

	mainWg.Wait()
	logger.L.Info("程序退出")
}

func handleLoop() {
	if !global.Aborting {
		checkLive()
		time.AfterFunc(time.Duration(checkDelay)*time.Second, handleLoop)
	}

}

func checkLive() {
	living, err := youtube.IsLiving(ytLiveUrl, ytCookie)
	if err != nil {
		logger.L.Errorf("检查直播间状态失败：%s", err.Error())
		return
	}
	if living {
		logger.L.Notice("频道直播中，开始推流")
		m3u8Url, err := youtube.GetLiveStreamURL(ytLiveUrl, ytCookie)
		if err != nil {
			logger.L.Errorf("获取直播流地址失败：%s", err.Error())
			return
		}

		mainWg.Add(1)
		ffParams := ffmpeg.GenReBoradcastParams(m3u8Url, pushURL, ytCookie)
		err = ffmpeg.StartCmd(ctx, ffmpegPath, ffParams, hideFFmpegLog)
		if err != nil {
			logger.L.Errorf("执行推流命令失败：%s", err.Error())
		}
		logger.L.Notice("推流结束")
		mainWg.Done()
	} else if !hideIdleInfo {
		logger.L.Info("频道未直播，等待...")
	}

}
