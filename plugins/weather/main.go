package plmm

import (
	"os"

	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
)

type Plmm struct {
	engine.PluginMagic
	Enable    bool   `yaml:"enable"`
	Dir       string `yaml:"dir"`
	Url       string `yaml:"url"`
	AppId     string `yaml:"appId"`
	AppSecret string `yaml:"appSecret"`
}

var (
	pluginInfo = &Plmm{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 /tianqi => 获取天气预报",
			Commands: []string{"/tianqi"},
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Plmm) OnRegister() {
	err := os.MkdirAll(plugin.RawConfig.Get("dir").(string), os.ModePerm)
	if err != nil {
		panic("init plmm img dir error: " + err.Error())
	}
}

func (p *Plmm) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() && msg.Content == pluginInfo.Commands[0] {
			// getWeather(msg)
		}
	}
}
