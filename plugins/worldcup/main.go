package worldcup

import (
	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
)

type WorldCup struct {
	engine.PluginMagic
}

var (
	keywords   = []string{"世界杯"}
	pluginInfo = &WorldCup{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 世界杯 => 获取2022世界杯赛程",
			Commands: keywords,
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *WorldCup) OnRegister() {
}

func (p *WorldCup) OnEvent(msg *robot.Message) {
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
