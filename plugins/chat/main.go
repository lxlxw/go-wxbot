package chat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/yqchilde/pkgs/log"

	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
)

type Chat struct {
	engine.PluginMagic
	Enable bool   `yaml:"enable"`
	Url    string `yaml:"url"`
}

var (
	keyword    = "@Bobo Bot"
	pluginInfo = &Chat{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 @Bot名称 => 跟它聊天",
			Commands: []string{keyword},
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Chat) OnRegister() {
}

func (p *Chat) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() && msg.IsAt() && strings.Contains(msg.Content, keyword) {
			getChatDetail(msg)
		}
	}
}

func getChatDetail(msg *robot.Message) {

	var chatConf Chat
	plugin.RawConfig.Unmarshal(&chatConf)

	chatMsg := strings.Trim(msg.Content, keyword)

	apiUrl := fmt.Sprintf("%s?key=free&appid=0&msg=%s", chatConf.Url, chatMsg)

	res, err := http.Get(apiUrl)
	if err != nil {
		log.Errorf("getChatDetail http get error: %v", err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getChatDetail read body error: %v", err)
		return
	}

	var resp ChatApiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getChatDetail unmarshal error: %v", err)
		return
	}
	if resp.Result != 0 {
		log.Errorf("getChatDetail api error: %v", resp.Content)
		return
	}

	msg.ReplyText(resp.Content)

}
