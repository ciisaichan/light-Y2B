package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ciisaichan/light-Y2B/common/setting"
	"github.com/ciisaichan/light-Y2B/utils"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ciisaichan/light-Y2B/ffmpeg"
	"github.com/ciisaichan/light-Y2B/global"
	"github.com/ciisaichan/light-Y2B/logger"
	"github.com/ciisaichan/light-Y2B/platform/youtube"
)

var (
	configPath string
	configName string

	ytLiveUrl string

	osChannel  chan os.Signal
	ctx        context.Context
	mainCancel context.CancelFunc
	mainWg     sync.WaitGroup
)

func init() {
	flag.StringVar(&configPath, "c-path", "./", "配置文件目录，不指定默认 当前目录")
	flag.StringVar(&configName, "c-name", "config", "配置文件名，不指定默认 config (文件后缀需要为.yaml)")
}

func main() {
	fmt.Printf("# Light-Y2B [%s]\n", global.Version)
	fmt.Println("# 基于 FFmpeg 的轻量级 Youtube 转播程序\n")
	flag.Parse()

	Setting, err := setting.NewSetting(configPath, configName)
	if err != nil {
		panic("读取配置文件失败：" + err.Error())
	}
	err = utils.ReadConfigToSetting(Setting)
	if err != nil {
		panic("识别配置文件失败：" + err.Error())
	}

	if global.Setting.Other.Proxy != "" {
		os.Setenv("http_proxy", global.Setting.Other.Proxy)
		os.Setenv("https_proxy", global.Setting.Other.Proxy)
	}

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

	if youtube.IsChannelURL(global.Setting.Live.YoutubeUrl) {
		ytLiveUrl = global.Setting.Live.YoutubeUrl + "/live"
	} else {
		ytLiveUrl = global.Setting.Live.YoutubeUrl
	}

	go handleLoop()
	logger.L.Info("开始监听直播间状态...")

	mainWg.Wait()
	logger.L.Info("程序退出")
}

func handleLoop() {
	if !global.Aborting {
		checkLive()
		time.AfterFunc(time.Duration(global.Setting.Other.CheckDelay)*time.Second, handleLoop)
	}

}

func checkLive() {
	living, err := youtube.IsLiving(ytLiveUrl, global.Setting.Live.YoutubeCookie)
	if err != nil {
		logger.L.Errorf("检查直播间状态失败：%s", err.Error())
		return
	}
	if living {
		logger.L.Notice("频道直播中，开始推流")
		m3u8Url, err := youtube.GetLiveStreamURL(ytLiveUrl, global.Setting.Live.YoutubeCookie)
		if err != nil {
			logger.L.Errorf("获取直播流地址失败：%s", err.Error())
			return
		}

		mainWg.Add(1)
		ffParams := ffmpeg.GenReBoradcastParams(m3u8Url, global.Setting.Live.BiliRtmpUrl, global.Setting.Live.YoutubeCookie)
		err = ffmpeg.StartCmd(ctx, global.Setting.Other.FFmpegPath, ffParams, global.Setting.Other.FFmpegLogs)
		if err != nil {
			logger.L.Errorf("执行推流命令失败：%s", err.Error())
		}
		logger.L.Notice("推流结束")
		mainWg.Done()
	} else if global.Setting.Other.IdleInfo {
		logger.L.Info("频道未直播，等待...")
	}

}
