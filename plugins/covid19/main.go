package covid19

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
)

type Covid19 struct {
	engine.PluginMagic
}

var (
	keyword    = "ç«æ"
	pluginInfo = &Covid19{
		PluginMagic: engine.PluginMagic{
			Desc:     "đ èŸć„ {city}ç«æ => è·ććźæ¶ç«æ",
			Commands: []string{keyword},
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Covid19) OnRegister() {
}

func (p *Covid19) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() && strings.Contains(msg.Content, pluginInfo.Commands[0]) {
			getC19(msg)
		}
	}
}

// result ç«ææ„èŻąç»æ
type checkresult struct {
	Data struct {
		Epidemic epidemic `json:"diseaseh5Shelf"`
	} `json:"data"`
}

// epidemic ç«ææ°æź
type epidemic struct {
	LastUpdateTime string  `json:"lastUpdateTime"`
	AreaTree       []*area `json:"areaTree"`
}

// area ććžç«ææ°æź
type area struct {
	Name   string `json:"name"`
	Adcode string `json:"adcode"`
	Today  struct {
		Confirm int         `json:"confirm"`
		Wzzadd  interface{} `json:"wzz_add"`
	} `json:"today"`
	Total struct {
		NowConfirm        int    `json:"nowConfirm"`
		Confirm           int    `json:"confirm"`
		Dead              int    `json:"dead"`
		Heal              int    `json:"heal"`
		HighRiskAreaNum   int    `json:"highRiskAreaNum"`
		MediumRiskAreaNum int    `json:"mediumRiskAreaNum"`
		Grade             string `json:"grade"`
		Wzz               int    `json:"wzz"`
	} `json:"total"`
	Children []*area `json:"children"`
}

func getC19(msg *robot.Message) {

	cityName := strings.Trim(msg.Content, keyword)

	data, time, err := queryEpidemic(cityName)
	if err != nil {
		return
	}
	if data == nil {
		return
	}

	var wzzadd string
	switch data.Today.Wzzadd.(type) {
	case string:
		op, _ := data.Today.Wzzadd.(string)
		wzzadd = op
	default:
		wzzadd = "0"
	}

	w_url := "https://news.qq.com/zt2020/page/feiyan.htm#/area?adcode=" + data.Adcode

	var str string
	str += "ă" + data.Name + "ăç«ææ°æź\n"
	str += "ä»æ„æ°ćąäșșæ°ïŒ" + strconv.Itoa(data.Today.Confirm) + "\n"
	str += "ä»æ„æ°ćąæ çç¶ïŒ" + wzzadd + "\n"
	//str += "ç°æçĄźèŻïŒ" + strconv.Itoa(data.Total.NowConfirm) + "\n"
	str += "çŽŻèźĄçĄźèŻäșșæ°ïŒ" + strconv.Itoa(data.Total.Confirm) + "\n"
	//str += "é«éŁé©ć°ćșïŒ" + strconv.Itoa(data.Total.HighRiskAreaNum) + "\n"
	//str += "äž­éŁé©ć°ćșïŒ" + strconv.Itoa(data.Total.MediumRiskAreaNum) + "\n"
	str += "æŽæ°æ¶éŽïŒ" + time + "\n\n"
	str += "èŻŠææ„çïŒ" + w_url

	msg.ReplyText(str)
}

func rcity(a *area, cityName string) *area {
	if a == nil {
		return nil
	}
	if a.Name == cityName {
		return a
	}
	for _, v := range a.Children {
		if v.Name == cityName {
			return v
		}
		c := rcity(v, cityName)
		if c != nil {
			return c
		}
	}
	return nil
}

func queryEpidemic(findCityName string) (citydata *area, times string, err error) {

	res, err := http.Get("https://api.inews.qq.com/newsqa/v1/query/inner/publish/modules/list?modules=statisGradeCityDetail,diseaseh5Shelf")
	if err != nil {
		return
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	var r checkresult
	err = json.Unmarshal(data, &r)
	if err != nil {
		return
	}
	citydata = rcity(r.Data.Epidemic.AreaTree[0], findCityName)
	return citydata, r.Data.Epidemic.LastUpdateTime, nil
}
