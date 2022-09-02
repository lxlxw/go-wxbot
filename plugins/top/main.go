package top

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
	"github.com/yqchilde/pkgs/log"
)

type Top struct {
	engine.PluginMagic
	Enable bool   `yaml:"enable"`
	Url    string `yaml:"url"`
	Key    string `yaml:"key"`
}

var (
	keywords   = []string{"知乎热榜", "知乎热搜", "B站热榜", "微信热榜", "微信热搜", "今日头条", "今日早报", "历史上的今天"}
	pluginInfo = &Top{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {知乎热榜} => 获取今日实时热榜",
			Commands: keywords,
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Top) OnRegister() {
}

func (p *Top) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() {
			for _, v := range keywords {
				if msg.Content == v {
					getTop(msg)
					return
				}
			}
		}
	}
}

var topMap = map[string]string{
	"知乎热榜":   "zhihu",
	"知乎热搜":   "zhihu",
	"B站热榜":   "bilibili",
	"微信热榜":   "weixin",
	"微信热搜":   "weixin",
	"今日头条":   "toutiao",
	"今日早报":   "toutiao",
	"历史上的今天": "hitory",
}

func getTop(msg *robot.Message) {

	var topConf Top
	plugin.RawConfig.Unmarshal(&topConf)

	apiUrl := fmt.Sprintf("%s?type=%s&token=%s", topConf.Url, topMap[msg.Content], topConf.Key)

	res, err := http.Get(apiUrl)
	if err != nil {
		log.Errorf("getTop http get error: %v", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getTop read body error: %v", err)
		return
	}

	var resp TopApiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getTop unmarshal error: %v", err)
		return
	}
	if resp.Code != 200 {
		log.Errorf("getTop api error: %v", resp.Msg)
		return
	}
	var str = ""
	str += "【" + resp.Result.Name + "】\n\n"

	for k, v := range resp.Result.List {
		if k >= 10 {
			continue
		}
		str += "标题：" + v.Title + "\n"
		str += "链接：" + v.Link + "\n"
		str += "\n\n"
	}

	msg.ReplyText(str)
}
