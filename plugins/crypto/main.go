package crypto

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
	keywords   = []string{"/btc", "/eth", "/比特币", "/以太坊"}
	otherName  = "数字货币"
	pluginInfo = &Crypto{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 数字货币{name} => 获取加密货币信息 || 示例：数字货币btc | /btc | /eth",
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
				if msg.Content == v || strings.Contains(msg.Content, otherName) {
					getCryptoDetail(msg, msg.Content)
					return
				}
			}
		}
	}
}

func getCryptoCode(keyword string) (string, string) {
	var code string
	var symbol string
	if keyword == "/比特币" || keyword == "/btc" {
		code = "BTCUSDT"
		symbol = "BTC"

	} else if keyword == "/以太坊" || keyword == "/eth" {
		code = "ETHUSDT"
		symbol = "ETH"
	} else if strings.Contains(keyword, otherName) {
		code = strings.Trim(keyword, otherName)
		symbol = code
		symbol = strings.ToUpper(symbol)
		code = code + "usdt"
		code = strings.ToUpper(code)
	}
	return code, symbol
}

func getCryptoDetail(msg *robot.Message, keyword string) {

	var cryptoConf Crypto
	plugin.RawConfig.Unmarshal(&cryptoConf)

	code, symbol := getCryptoCode(keyword)

	apiUrl := fmt.Sprintf("%s?symbol=%s", cryptoConf.Url, code)
	log.Println(apiUrl)
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
	priceFloat, _ := strconv.ParseFloat(resp.Price, 64)
	price := fmt.Sprintf("%.3f", priceFloat)

	detail := fmt.Sprintf(`%s Price：$%s`, symbol, price)

	msg.ReplyText(detail)
}
