package crazykfc

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
)

type CrazyKFC struct{ engine.PluginMagic }

var (
	pluginInfo = &CrazyKFC{
		engine.PluginMagic{
			Desc:     "ð è¾å¥ /kfc => è·åè¯å¾·åºç¯çææå",
			Commands: []string{"/kfc"},
		},
	}
	_        = engine.InstallPlugin(pluginInfo)
	sentence []string
)

func (p *CrazyKFC) OnRegister() {
	resp, err := getCrazyKFCSentence()
	if err != nil {
		return
	}
	for i := range resp {
		sentence = append(sentence, resp[i].Text)
	}
}

func (p *CrazyKFC) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() && strings.HasPrefix(msg.Content, pluginInfo.Commands[0]) {
			if len(sentence) > 0 {
				rand.Seed(time.Now().UnixNano())
				msg.ReplyText(sentence[rand.Intn(len(sentence))])
			} else {
				msg.ReplyText("æ¥è¯¢å¤±è´¥ï¼è¿ä¸å®ä¸æ¯bugð¤")
			}
		}
	}
}

type apiResponse struct {
	Index int    `json:"index"`
	Text  string `json:"text"`
}

func getCrazyKFCSentence() ([]apiResponse, error) {
	api := "https://fastly.jsdelivr.net/gh/Nthily/KFC-Crazy-Thursday@main/kfc.json"
	resp, err := http.Get(api)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	readAll, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data []apiResponse
	if err := json.Unmarshal(readAll, &data); err != nil {
		return nil, err
	}
	return data, nil
}
