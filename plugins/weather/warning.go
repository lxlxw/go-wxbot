package weather

import (
	"encoding/json"
	"fmt"
	"github.com/yqchilde/pkgs/log"
	"io"
	"net/http"
	"strings"
)

type WarningInfo struct{}

func (c *WarningInfo) GetInfo(url string, locationID string) string {

	detail := "灾害预警：无"
	url = fmt.Sprintf("%s?location=%s&key=%s", url, locationID, weatherConf.Key)

	res, err := http.Get(url)
	if err != nil {
		log.Errorf("getWeatherWarningInfo http get error: %v", err)
		return detail
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getWeatherWarningInfo read body error: %v", err)
		return detail
	}

	var resp ResponseWeatherWarning
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getWeatherWarningInfo unmarshal error: %v", err)
		return detail
	}
	if resp.Code != "200" {
		log.Errorf("getWeatherWarningInfo api error: code not 200" + string(url))
		return detail
	}

	if len(resp.Warning) > 0 {
		detail = strings.Replace(detail, "无", resp.Warning[0].Text, 1)
	}

	return detail
}
