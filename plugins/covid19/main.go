package covid19

import (
	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
)

type Covid19 struct {
	engine.PluginMagic
}

var (
	keywords   = []string{"疫情"}
	pluginInfo = &Covid19{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {city}疫情 => 获取疫情情况",
			Commands: keywords,
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Covid19) OnRegister() {
}

func (p *Covid19) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() {
			for _, v := range keywords {
				if msg.Content == v {
					//getJx3Detail(msg, msg.Content)
					return
				}
			}
		}
	}
}
