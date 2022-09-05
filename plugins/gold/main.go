package gold

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

type Gold struct {
	engine.PluginMagic
	Enable bool   `yaml:"enable"`
	Url    string `yaml:"url"`
	Key    string `yaml:"key"`
}

var (
	keyword    = "黄金价格"
	pluginInfo = &Gold{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {黄金价格} => 获取黄金实时价格",
			Commands: []string{keyword},
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Gold) OnRegister() {
}

func (p *Gold) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() && strings.Contains(msg.Content, pluginInfo.Commands[0]) {
			GetGold(msg)
		}
	}
}

func GetGold(msg *robot.Message) {
	var goldConf Gold
	plugin.RawConfig.Unmarshal(&goldConf)

	apiUrl := fmt.Sprintf("%s?key=%s", goldConf.Url, goldConf.Key)

	res, err := http.Get(apiUrl)
	if err != nil {
		log.Errorf("getOil http get error: %v", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getGold read body error: %v", err)
		return
	}

	var resp ResponseGoldApi
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getGold unmarshal error: %v", err)
		return
	}
	if resp.ErrorCode != 0 {
		log.Errorf("getGold api error: %v", resp.Reason)
		return
	}

	str := "【黄金价格】\n"
	for _, v := range resp.Result {
		for _, v := range v {
			str += "品种：" + v.Variety + "\n"
			str += "最新价：" + v.Latestpri + "\n"
			str += "开盘价：" + v.Openpri + "\n"
			str += "最高价：" + v.Maxpri + "\n"
			str += "最低价：" + v.Minpri + "\n"
			str += "涨跌幅：" + v.Limit + "\n"
			str += "昨收价：" + v.Yespri + "\n"
			str += "总成交量：" + v.Totalvol + "\n"
			str += "更新时间：" + v.Time + "\n"
		}
	}

	_, err = msg.ReplyText(str)
	if err != nil {
		println()
		return
	}
}
