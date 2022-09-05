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
	AppId     string `yaml:"appId"`
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

func getCryptoCode(content string) (string, string) {
	var code string
	var symbol string
	if strings.Contains(content, keyword) {
		code = strings.Trim(content, keyword)
		symbol = code
		symbol = strings.ToUpper(symbol)
		code = code + "-USD"
		code = strings.ToUpper(code)
	}
	return code, symbol
}

func getCryptoDetail(msg *robot.Message) {

	code, symbol := getCryptoCode(msg.Content)

	apiUrl := fmt.Sprintf("https://api.blockchain.com/v3/exchange/tickers/%s", code)

	res, err := http.Get(apiUrl)
	if err != nil {
		log.Errorf("getCryptoDetail http get error: %v", err)
		return
	}
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
	price := fmt.Sprintf("%.4f", resp.Last_trade_price)

	detail := fmt.Sprintf(`%s Price：$%s`, symbol, price)

	msg.ReplyText(detail)
}
