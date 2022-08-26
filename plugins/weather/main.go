package weather

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

type Weather struct {
	engine.PluginMagic
	Enable    bool   `yaml:"enable"`
	Url       string `yaml:"url"`
	AppId     string `yaml:"appId"`
	AppSecret string `yaml:"appSecret"`
}

var (
	keyword    = "天气"
	pluginInfo = &Weather{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {city}天气 => 获取天气预报",
			Commands: []string{keyword},
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Weather) OnRegister() {
}

func (p *Weather) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() && strings.Contains(msg.Content, pluginInfo.Commands[0]) {
			getWeatherDetail(msg)
		}
	}
}

func getWeatherDetail(msg *robot.Message) {

	var weatherConf Weather
	plugin.RawConfig.Unmarshal(&weatherConf)

	cityName := strings.Trim(msg.Content, keyword)

	apiUrl := fmt.Sprintf("%s/%s?app_id=%s&app_secret=%s", weatherConf.Url, cityName, weatherConf.AppId, weatherConf.AppSecret)

	res, err := http.Get(apiUrl)
	if err != nil {
		log.Errorf("getWeatherDetail http get error: %v", err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getWeatherDetail read body error: %v", err)
		return
	}

	var resp WeatherApiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getWeatherDetail unmarshal error: %v", err)
		return
	}
	if resp.Code != 1 {
		log.Errorf("getWeatherDetail api error: %v", resp.Msg)
		return
	}

	detail := fmt.Sprintf(`%s今天天气，温度为 %s ，天气 %s，%s %s，相对湿度 %s`, resp.Data.Address,
		resp.Data.Temp, resp.Data.Weather, resp.Data.WindDirection, resp.Data.WindPower, resp.Data.Humidity)

	msg.ReplyText(detail)

}
