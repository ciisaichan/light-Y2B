package setting

import (
	"github.com/spf13/viper"
)

type Setting struct {
	vp *viper.Viper
}

func NewSetting(configPath string, configName string) (*Setting, error) {
	Vp := viper.New()
	Vp.SetConfigName(configName) // 配置文件名
	Vp.AddConfigPath(configPath) // 配置文件路径
	Vp.SetConfigType("yaml")     // 配置文件类型
	err := Vp.ReadInConfig()     // 读取配置文件
	if err != nil {
		return nil, err
	}
	return &Setting{Vp}, nil
}
