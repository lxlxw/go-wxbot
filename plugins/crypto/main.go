package crypto

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

type Crypto struct {
	engine.PluginMagic
	Enable    bool   `yaml:"enable"`
	Url       string `yaml:"url"`
	AppSecret string `yaml:"appSecret"`
}

var (
	keywords   = []string{"$btc", "$eth"}
	keyword    = "$"
	pluginInfo = &Crypto{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 $btc => 获取加密货币信息 || 示例：$btc | $eth",
			Commands: keywords,
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Crypto) OnRegister() {
}

func (p *Crypto) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() {
			for _, v := range keywords {
				if msg.Content == v || strings.Contains(msg.Content, keyword) {
					getCryptoDetail(msg)
					return
				}
			}
		}
	}
}

func getCryptoDetail(msg *robot.Message) {

	var cryptoConf Crypto
	plugin.RawConfig.Unmarshal(&cryptoConf)

	code := strings.Trim(msg.Content, keyword)
	code = strings.ToUpper(code)

	apiUrl := fmt.Sprintf("%s?symbol=%s&interval=0", cryptoConf.Url, code)

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)

	req.Header.Add("coinglassSecret", cryptoConf.AppSecret)
	if err != nil {
		return
	}
	res, _ := client.Do(req)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getCryptoDetail read body error: %v", err)
		return
	}

	var resp CryptoApiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getCryptoDetail unmarshal error: %v", err)
		return
	}
	if resp.Code != "0" {
		return
	}
	var str string
	if len(resp.List) <= 0 {
		return
	}

	for _, v := range resp.List {
		if v.ExchangeName != "All" {
			price := fmt.Sprintf("%f", v.Price)
			str += code + " Price：" + price + "\n\n"
			str += "数据来源：" + v.ExchangeName
			break
		}
	}

	msg.ReplyText(str)
}
