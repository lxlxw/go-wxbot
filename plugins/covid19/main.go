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
	keyword    = "疫情"
	pluginInfo = &Covid19{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {city}疫情 => 获取实时疫情",
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

// result 疫情查询结果
type checkresult struct {
	Data struct {
		Epidemic epidemic `json:"diseaseh5Shelf"`
	} `json:"data"`
}

// epidemic 疫情数据
type epidemic struct {
	LastUpdateTime string  `json:"lastUpdateTime"`
	AreaTree       []*area `json:"areaTree"`
}

// area 城市疫情数据
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
	str += "【" + data.Name + "】疫情数据\n"
	str += "今日新增人数：" + strconv.Itoa(data.Today.Confirm) + "\n"
	str += "今日新增无症状：" + wzzadd + "\n"
	//str += "现有确诊：" + strconv.Itoa(data.Total.NowConfirm) + "\n"
	str += "累计确诊人数：" + strconv.Itoa(data.Total.Confirm) + "\n"
	//str += "高风险地区：" + strconv.Itoa(data.Total.HighRiskAreaNum) + "\n"
	//str += "中风险地区：" + strconv.Itoa(data.Total.MediumRiskAreaNum) + "\n"
	str += "更新时间：" + time + "\n\n"
	str += "详情查看：" + w_url

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
