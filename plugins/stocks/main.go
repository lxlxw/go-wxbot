package stocks

import (
	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
)

type Stocks struct {
	engine.PluginMagic
}

var (
	keywords   = []string{"股票"}
	pluginInfo = &Stocks{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {股票名} => 获取股票情况 | 示例：特斯拉",
			Commands: keywords,
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Stocks) OnRegister() {
}

func (p *Stocks) OnEvent(msg *robot.Message) {
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
