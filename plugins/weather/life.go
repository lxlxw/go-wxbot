package weather

import (
	"encoding/json"
	"fmt"
	"github.com/yqchilde/pkgs/log"
	"io"
	"net/http"
)

type LifeInfo struct{}

func (c *LifeInfo) GetInfo(locationID string) string {

	detail := ""
	url := fmt.Sprintf("%s?location=%s&key=%s&type=1,3,4,5", weatherConf.LifeUrl, locationID, weatherConf.Key)

	res, err := http.Get(url)
	if err != nil {
		log.Errorf("getWeatherLifeInfo http get error: %v", err)
		return detail
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getWeatherLifeInfo read body error: %v", err)
		return detail
	}

	var resp ResponseWeatherLife
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getWeatherLifeInfo unmarshal error: %v", err)
		return detail
	}
	if resp.Code != "200" {
		log.Errorf("getWeatherLifeInfo api error: code not 200" + string(url))
		return detail
	}

	for _, v := range resp.Daily {
		detail += v.Name + "：" + v.Text + "\n"
	}

	return detail
}
