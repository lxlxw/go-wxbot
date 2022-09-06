package weather

import (
	"encoding/json"
	"fmt"
	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
	"github.com/yqchilde/pkgs/log"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Weather struct {
	engine.PluginMagic
	Enable      bool   `yaml:"enable"`
	Key         string `yaml:"key"`
	LocationUrl string `yaml:"locationUrl"`
	ActualUrl   string `yaml:"actualUrl"`
	LifeUrl     string `yaml:"lifeUrl"`
	AirUrl      string `yaml:"airUrl"`
	WarningUrl  string `yaml:"warningUrl"`
	SunUrl      string `yaml:"sunUrl"`
}

var (
	keyword    = "天气"
	pluginInfo = &Weather{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {city}天气 => 获取天气预报 || 示例：北京天气",
			Commands: []string{keyword},
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

var weatherConf Weather

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

	var urlMap = map[string]string{
		"实时天气": weatherConf.ActualUrl,
		"空气指数": weatherConf.AirUrl,
		"日出日落": weatherConf.SunUrl,
		"生活指数": weatherConf.LifeUrl,
		"气象预警": weatherConf.WarningUrl,
	}
	detail := ""
	plugin.RawConfig.Unmarshal(&weatherConf)

	cityName := strings.Trim(msg.Content, keyword)

	localtionUrl := fmt.Sprintf("%s?location=%s&key=%s", weatherConf.LocationUrl, url.QueryEscape(cityName), weatherConf.Key)

	res, err := http.Get(localtionUrl)
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

	var resp ResponseWeatherLocaltion
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getWeatherDetail unmarshal error: %v", localtionUrl)
		return
	}
	if resp.Code != "200" {
		log.Errorf("getWeatherDetail api error: 未查询到城市信息")
		detail = "未查询到城市信息"
		msg.ReplyText(detail)
		return
	}

	locationID := resp.Location[0].ID

	detalMap := map[string]string{}

	for k, v := range urlMap {
		info, err := Factory(k)
		if err != nil {
			log.Errorf("error: %v", err)
			return
		}
		detalMap[k] = info.GetInfo(v, locationID)
	}

	//TODO 并发请求
	//wg := &sync.WaitGroup{}
	//for k, v := range urlMap {
	//	wg.Add(1)
	//	go func(url string) {
	//		info, err := Factory(k)
	//		if err != nil {
	//			log.Errorf("error: %v", err)
	//			return
	//		}
	//		detalMap[k] = info.GetInfo(url, locationID)
	//	}(v)
	//	wg.Done()
	//}
	//wg.Wait()

	detail = resp.Location[0].Name + "天气\n"
	for k, _ := range urlMap {
		detail += detalMap[k]
	}

	msg.ReplyText(detail)

}
