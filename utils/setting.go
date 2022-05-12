package utils

import (
	"github.com/ciisaichan/light-Y2B/common/setting"
	"github.com/ciisaichan/light-Y2B/global"
)

// 读取配置文件
func ReadConfigToSetting(setting *setting.Setting) error {
	var err error
	err = setting.ReadSection(&global.Setting)
	if err != nil {
		return err
	}
	return nil
}
