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

var infoMap = []string{"实时天气", "空气指数", "日出日落", "生活指数", "气象预警"}

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

	for _, v := range infoMap {
		info, err := Factory(v)
		if err != nil {
			log.Errorf("error: %v", err)
			return
		}
		detalMap[v] = info.GetInfo(locationID)
	}

	//TODO 并发请求
	//wg := &sync.WaitGroup{}
	//wg.Add(len(urlMap))
	//for k, v := range urlMap {
	//	go func(k string, url string) {
	//		defer wg.Done()
	//		info, err := Factory(k)
	//		if err != nil {
	//			log.Errorf("error: %v", err)
	//			return
	//		}
	//		detalMap[k] = info.GetInfo(url, locationID)
	//	}(k, v)
	//}
	//wg.Wait()

	detail = resp.Location[0].Name + "天气\n"
	for _, v := range infoMap {
		detail += detalMap[v]
	}

	msg.ReplyText(detail)

}
