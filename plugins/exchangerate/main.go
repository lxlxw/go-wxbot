package exchangerate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
	"github.com/yqchilde/pkgs/log"
)

type ExchangeRate struct {
	engine.PluginMagic
	Enable    bool   `yaml:"enable"`
	Url       string `yaml:"url"`
	AppId     string `yaml:"appId"`
	AppSecret string `yaml:"appSecret"`
}

var curr = map[string]string{
	"人民币":   "CNY",
	"美元":    "USD",
	"欧元":    "EUR",
	"日元":    "JPY",
	"港币":    "HKD",
	"英镑":    "GBP",
	"澳元":    "AUD",
	"新西兰元":  "NZD",
	"新加坡元":  "SGD",
	"瑞士法郎":  "CHF",
	"加元":    "CAD",
	"俄罗斯卢布": "RUB",
	"卢布":    "RUB",
	"韩元":    "KRW",
	"丹麦克朗":  "DKK",
	"瑞典克朗":  "SEK",
	"泰铢":    "THB",
	"加拿大元":  "CAD",
	"澳门元":   "MOP",
	"墨西哥元":  "MXN",
	"泰元":    "THB",
	"台币":    "TWD",
	"越南盾":   "VND",
}

var (
	keyword    = "兑"
	pluginInfo = &ExchangeRate{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {美元}兑{人民币} => 获取汇率 | 示例：美元兑人民币",
			Commands: []string{keyword},
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *ExchangeRate) OnRegister() {
}

func (p *ExchangeRate) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() && strings.Contains(msg.Content, pluginInfo.Commands[0]) {
			getExDetail(msg)
		}
	}
}
func getExDetail(msg *robot.Message) {
	var exConf ExchangeRate
	plugin.RawConfig.Unmarshal(&exConf)

	curMsg := strings.Split(msg.Content, keyword)
	if len(curMsg) != 2 {
		return
	}

	var from string
	var to string
	for k, v := range curMsg {
		if _, ok := curr[v]; ok {
			if k == 0 {
				from = curr[v]
			} else if k == 1 {
				to = curr[v]
			}
		} else {
			return
		}
	}
	apiUrl := fmt.Sprintf("%s?key=%s&from=%s&to=%s", exConf.Url, exConf.AppSecret, from, to)

	res, err := http.Get(apiUrl)
	if err != nil {
		log.Errorf("getExDetail http get error: %v", err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getExDetail read body error: %v", err)
		return
	}

	var resp ExApiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getExDetail unmarshal error: %v", err)
		return
	}
	if resp.Code != 0 {
		log.Errorf("getExDetail api error: %v", resp.Msg)
		return
	}
	var str string
	str += resp.Data[0].CurrencyF_Name + "(" + resp.Data[0].CurrencyF + ")/" + resp.Data[0].CurrencyT_Name + "(" + resp.Data[0].CurrencyT + ")\n"
	str += "当前汇率为：" + resp.Data[0].Exchange + "\n\n"

	str += resp.Data[1].CurrencyF_Name + "(" + resp.Data[1].CurrencyF + ")/" + resp.Data[1].CurrencyT_Name + "(" + resp.Data[1].CurrencyT + ")\n"
	str += "当前汇率为：" + resp.Data[1].Exchange + "\n\n"

	str += "更新时间：" + resp.Data[1].UpdateTime + "\n"

	msg.ReplyText(str)
}
