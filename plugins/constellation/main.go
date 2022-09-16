package constellation

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
	"github.com/yqchilde/pkgs/log"
)

type Constellation struct {
	engine.PluginMagic
	Enable bool   `yaml:"enable"`
	Url    string `yaml:"url"`
	Key    string `yaml:"key"`
}

var dateName = map[string]string{
	"今日": "today",
	"明日": "tomorrow",
	"本周": "week",
	"本月": "month",
	"今年": "year",
}

var dateKeyword = []string{"今日", "明日", "本周", "本月", "今年"}

var (
	keyword    = "运势"
	pluginInfo = &Constellation{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {星座}{日期}运势 => 获取星座运势 | 示例：白羊座今日运势",
			Commands: []string{keyword},
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Constellation) OnRegister() {
}

func (p *Constellation) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() && strings.Contains(msg.Content, keyword) {
			geConstellationDetail(msg)
		}
		if msg.Content == "星座运势" {
			str := "🚀 输入 {星座}{日期}运势 => 获取星座运势\n\n"
			str += "示例：\n"
			str += "1 - 白羊座今日运势 \n"
			str += "2 - 双鱼座明日运势 \n"
			str += "3 - 天蝎座本周运势 \n"
			str += "4 - 金牛座本月运势 \n"
			str += "5 - 金牛座今年运势"
			msg.ReplyText(str)
		}
	}
}

func geConstellationDetail(msg *robot.Message) {
	var exConf Constellation
	plugin.RawConfig.Unmarshal(&exConf)

	consName := msg.Content
	var dateType string
	consName = strings.Trim(consName, keyword)
	for _, v := range dateKeyword {
		if strings.Contains(consName, v) {
			dateType = dateName[v]
			consName = strings.Trim(consName, v)
			break
		}
	}

	data := url.Values{}
	data.Add("consName", consName)
	data.Add("type", dateType)
	data.Add("key", exConf.Key)

	u, _ := url.ParseRequestURI("https://web.juhe.cn")

	u.Path = "/constellation/getAll"
	u.RawQuery = data.Encode()
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return
	}
	res, _ := client.Do(req)

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("geConstellationDetail read body error: %v", err)
		return
	}
	var str string
	if dateType == "today" || dateType == "tomorrow" {
		str = geConsTodayDetail(body)
		msg.ReplyText(str)
	} else if dateType == "week" || dateType == "month" {
		str = geConsWeekDetail(body)
		msg.ReplyText(str)
	} else if dateType == "year" {
		str = geConsYearDetail(body)
		msg.ReplyText(str)
	} else {
		return
	}
}

func geConsTodayDetail(body []byte) string {
	var str string
	var resp ConTodayApiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Println(err)
		return ""
	}
	if resp.Code != 0 {
		return ""
	}
	str += "【" + resp.Datetime + resp.Name + "星座运势】" + "\n\n"
	str += "今日概述：" + resp.Summary + "\n"
	str += "速配星座：" + resp.QFriend + "\n"
	str += "幸运色：" + resp.Color + "\n"
	str += "幸运数字：" + strconv.Itoa(resp.Number) + "\n"
	str += "健康指数：" + resp.Health + "\n"
	str += "爱情指数：" + resp.Love + "\n"
	str += "财运指数：" + resp.Money + "\n"
	str += "工作指数：" + resp.Work + "\n"
	str += "综合指数：" + resp.All
	return str
}

func geConsWeekDetail(body []byte) string {
	var str string
	var resp CoWeekApiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Println(err)
		return ""
	}
	if resp.Code != 0 {
		return ""
	}
	str += "【" + resp.Date + resp.Name + "星座运势】" + "\n\n"

	str += "1、健康：" + resp.Health + "\n"
	str += "2、爱情：" + resp.Love + "\n"
	str += "3、财运：" + resp.Money + "\n"
	str += "4、工作：" + resp.Work + "\n"
	if resp.All != "" {
		str += "5、综合指数：" + resp.All
	}
	return str
}

func geConsYearDetail(body []byte) string {
	var str string
	var resp CoYearApiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Println(err)
		return ""
	}
	if resp.Code != 0 {
		return ""
	}
	str += "【" + resp.Date + resp.Name + "星座运势】" + "\n\n"

	str += "1、年度密码：" + "\n"
	str += "<概述>：" + resp.Mima.Info + "\n"
	str += "<说明>：" + resp.Mima.Text[0] + "\n\n"
	str += "2、事业运：" + resp.Career[0] + "\n"
	str += "3、感情运：" + resp.Love[0] + "\n"
	str += "4、财运：" + resp.Finance[0] + "\n"
	return str
}
