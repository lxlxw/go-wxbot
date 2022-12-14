package worldcup

import (
	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
)

type WorldCup struct {
	engine.PluginMagic
}

var (
	keywords   = []string{"δΈηζ―"}
	pluginInfo = &WorldCup{
		PluginMagic: engine.PluginMagic{
			Desc:     "π θΎε₯ δΈηζ― => θ·ε2022δΈηζ―θ΅η¨",
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
