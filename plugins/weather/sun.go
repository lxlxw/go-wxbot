package weather

import (
	"encoding/json"
	"fmt"
	"github.com/yqchilde/pkgs/log"
	"io"
	"net/http"
	"time"
)

type SunInfo struct{}

func (c *SunInfo) GetInfo(locationID string) string {

	detail := ""
	url := fmt.Sprintf("%s?location=%s&key=%s&date=%s", weatherConf.SunUrl, locationID, weatherConf.Key, time.Now().Format("20060102"))

	res, err := http.Get(url)
	if err != nil {
		log.Errorf("getWeatherSunInfo http get error: %v", err)
		return detail
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getWeatherSunInfo read body error: %v", err)
		return detail
	}

	var resp ResponseWeatherSun
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getWeatherSunInfo unmarshal error: %v", err)
		return detail
	}
	if resp.Code != "200" {
		log.Errorf("getWeatherActualInfo api error: code not 200" + string(url))
		return detail
	}
	var sunrise, sunset string
	if len(resp.Sunrise) > 11 {
		sunrise = resp.Sunrise[11 : len(resp.Sunrise)-6]
	} else {
		sunrise = resp.Sunrise
	}

	if len(resp.Sunset) > 11 {
		sunset = resp.Sunset[11 : len(resp.Sunset)-6]
	} else {
		sunset = resp.Sunset
	}

	detail = fmt.Sprintf("日出：%s\n日落：%s\n",
		sunrise,
		sunset)

	return detail
}
