package oil

import (
	"encoding/json"
	"fmt"
	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
	"github.com/yqchilde/pkgs/log"
	"io"
	"net/http"
	"strings"
)

type Oil struct {
	engine.PluginMagic
	Enable bool   `yaml:"enable"`
	Url    string `yaml:"url"`
	Key    string `yaml:"key"`
}

var (
	keyword    = "油价"
	pluginInfo = &Oil{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {province}油价 => 获取实时油价 || 示例：福建油价",
			Commands: []string{keyword},
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Oil) OnRegister() {
}

func (p *Oil) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() && strings.Contains(msg.Content, pluginInfo.Commands[0]) {
			getOil(msg)
		}
	}
}

func getOil(msg *robot.Message) {
	var oidConf Oil
	plugin.RawConfig.Unmarshal(&oidConf)

	provinceName := strings.Trim(msg.Content, keyword)

	apiUrl := fmt.Sprintf("%s?key=%s", oidConf.Url, oidConf.Key)

	res, err := http.Get(apiUrl)
	if err != nil {
		log.Errorf("getOil http get error: %v", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getOil read body error: %v", err)
		return
	}

	var resp ResponseOilApi
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getOil unmarshal error: %v", err)
		return
	}
	if resp.ErrorCode != 0 {
		log.Errorf("getOil api error: %v", resp.Reason)
		return
	}

	detail := "未查询到该城市油价"
	for _, v := range resp.Result {
		if v.City == provinceName {
			detail = fmt.Sprintf("【%s今天油价】\n92汽油价格：%s\n95汽油价格：%s\n98汽油价格：%s\n0号柴油价格：%s\n", v.City,
				v.Oil92h, v.Oil95h, v.Oil98h, v.Oil0h)
		}
	}

	_, err = msg.ReplyText(detail)
	if err != nil {
		return
	}

}
