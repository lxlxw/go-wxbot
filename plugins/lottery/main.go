package lottery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
	"github.com/yqchilde/pkgs/log"
)

type Lottery struct {
	engine.PluginMagic
	Enable    bool   `yaml:"enable"`
	Url       string `yaml:"url"`
	AppId     string `yaml:"appId"`
	AppSecret string `yaml:"appSecret"`
}

var (
	keywords   = []string{"双色球", "超级大乐透", "大乐透", "七乐彩", "福彩3D", "七星彩", "ssq", "qlc", "fc", "dlt", "qxc"}
	pluginInfo = &Lottery{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 双色球 => 获取彩票信息 || 示例：双色球|大乐透",
			Commands: keywords,
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Lottery) OnRegister() {
}

func (p *Lottery) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() {
			for _, v := range keywords {
				if msg.Content == v {
					getLotteryDetail(msg, v)
					return
				}
			}
		}
	}
}

func getLotteryCode(keyword string) string {
	var code string
	if keyword == "双色球" || keyword == "ssq" {
		code = "ssq"
	} else if keyword == "超级大乐透" || keyword == "大乐透" || keyword == "dlt" {
		code = "cjdlt"
	} else if keyword == "七乐彩" || keyword == "qlc" {
		code = "qlc"
	} else if keyword == "七星彩" || keyword == "qxc" {
		code = "qxc"
	} else if keyword == "福彩3D" || keyword == "fc3d" {
		code = "fc3d"
	}
	return code
}

func getLotteryDetail(msg *robot.Message, keyword string) {

	var lotteryConf Lottery
	plugin.RawConfig.Unmarshal(&lotteryConf)

	code := getLotteryCode(keyword)

	apiUrl := fmt.Sprintf("%s?code=%s&app_id=%s&app_secret=%s", lotteryConf.Url, code, lotteryConf.AppId, lotteryConf.AppSecret)
	res, err := http.Get(apiUrl)
	if err != nil {
		log.Errorf("getLotteryDetail http get error: %v", err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getLotteryDetail read body error: %v", err)
		return
	}

	var resp LotteryApiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getLotteryDetail unmarshal error: %v", err)
		return
	}
	if resp.Code != 1 {
		log.Errorf("getLotteryDetail api error: %v", resp.Msg)
		return
	}

	detail := fmt.Sprintf(`第%s期%s开奖结果：`+"\n"+"--------------"+"\n"+`开奖号码：`+"\n"+`%s`+"\n"+"--------------"+"\n"+`开奖日期：`+"\n"+`%s`, resp.Data.Expect, resp.Data.Name,
		resp.Data.OpenCode, resp.Data.Time)

	msg.ReplyText(detail)
}
