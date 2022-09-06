package weather

import (
	"encoding/json"
	"fmt"
	"github.com/yqchilde/pkgs/log"
	"io"
	"net/http"
)

type ActualInfo struct{}

func (c *ActualInfo) GetInfo(url string, locationID string) string {

	detail := ""
	url = fmt.Sprintf("%s?location=%s&key=%s", url, locationID, weatherConf.Key)

	res, err := http.Get(url)
	if err != nil {
		log.Errorf("getWeatherActualInfo http get error: %v", err)
		return detail
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getWeatherActualInfo read body error: %v", err)
		return detail
	}

	var resp ResponseWeatherActual
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getWeatherActualInfo unmarshal error: %v", err)
		return detail
	}
	if resp.Code != "200" {
		log.Errorf("getWeatherActualInfo api error: code not 200")
		return detail
	}

	detail = fmt.Sprintf("温度：%s℃\n体感温度：%s℃\n气象：%s\n风向：%s\n风力等级：%s\n风速：%s\n相对湿度：%s%%\n能见度：%s公里\n",
		resp.Now.Temp,
		resp.Now.Feelslike,
		resp.Now.Text,
		resp.Now.Winddir,
		resp.Now.Windscale,
		resp.Now.Windspeed,
		resp.Now.Humidity,
		resp.Now.Vis)

	return detail
}
