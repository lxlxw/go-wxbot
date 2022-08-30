package oil

import (
	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
)

type Oil struct {
	engine.PluginMagic
}

var (
	keywords   = []string{"油价"}
	pluginInfo = &Oil{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {city}油价 => 获取实时油价 | 示例：北京油价",
			Commands: keywords,
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Oil) OnRegister() {
}

func (p *Oil) OnEvent(msg *robot.Message) {
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
