package weather

import (
	"encoding/json"
	"fmt"
	"github.com/yqchilde/pkgs/log"
	"io"
	"net/http"
)

type AirInfo struct{}

func (c *AirInfo) GetInfo(locationID string) string {

	detail := ""
	url := fmt.Sprintf("%s?location=%s&key=%s", weatherConf.AirUrl, locationID, weatherConf.Key)

	res, err := http.Get(url)
	if err != nil {
		log.Errorf("getWeatherAirInfo http get error: %v", err)
		return detail
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getWeatherAirInfo read body error: %v", err)
		return detail
	}

	var resp ResponseWeatherAir
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getWeatherAirInfo unmarshal error: %v", err)
		return detail
	}
	if resp.Code != "200" {
		log.Errorf("getWeatherAirInfo api error: code not 200")
		return detail
	}

	detail = fmt.Sprintf("空气质量指数：%s\nPM2.5：%s\n",
		resp.Now.Aqi,
		resp.Now.Pm2P5)

	return detail
}
