package gold

import (
	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
)

type Gold struct {
	engine.PluginMagic
}

var (
	keywords   = []string{"黄金价格"}
	pluginInfo = &Gold{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {黄金价格} => 获取黄金实时价格",
			Commands: keywords,
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Gold) OnRegister() {
}

func (p *Gold) OnEvent(msg *robot.Message) {
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
