package bootstrap

import (
	"log"

	"github.com/ciisaichan/light-Y2B/common/setting"
	"github.com/ciisaichan/light-Y2B/global"
)

func init() {
	// 读取配置文件
	{
		Setting, err := setting.NewSetting()
		if err != nil {
			log.Println("init.setting.NewSetting():", err)
		}
		err = ReadConfigToSetting(Setting)
		if err != nil {
			log.Fatalf("init.setupSetting err: %v", err)
		}
	}
}

// 读取配置文件
func ReadConfigToSetting(setting *setting.Setting) error {
	var err error
	err = setting.ReadSection("Live", &global.LiveSetting)
	if err != nil {
		return err
	}
	return nil
}
