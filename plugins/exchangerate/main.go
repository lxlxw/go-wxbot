package exchangerate

import (
	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
)

type ExchangeRate struct {
	engine.PluginMagic
}

var (
	keywords   = []string{"汇率"}
	pluginInfo = &ExchangeRate{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {美元}兑{人民币} => 获取汇率 | 示例：美元兑人名币",
			Commands: keywords,
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *ExchangeRate) OnRegister() {
}

func (p *ExchangeRate) OnEvent(msg *robot.Message) {
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
