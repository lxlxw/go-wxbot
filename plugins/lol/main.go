package lol

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
)

type Lol struct {
	engine.PluginMagic
}

var month = map[string]string{
	"January":   "1",
	"February":  "2",
	"March":     "3",
	"April":     "4",
	"May":       "5",
	"June":      "6",
	"July":      "7",
	"August":    "8",
	"September": "9",
	"October":   "10",
	"November":  "11",
	"December":  "12",
}

var (
	keyword    = "lol赛程"
	pluginInfo = &Lol{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 lol赛程 => 获取最近LOL赛程",
			Commands: []string{keyword},
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Lol) OnRegister() {
}

func (p *Lol) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() && msg.Content == keyword {
			getLOL(msg)
		}
	}
}

func getLOL(msg *robot.Message) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://lol.fandom.com/wiki/League_of_Legends_Esports_Wiki", nil)

	if err != nil {
		return
	}
	res, _ := client.Do(req)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	var str = ""
	var rankInt = 1
	doc.Find(".frontpage-upcoming-matches tbody tr").Each(func(i int, s *goquery.Selection) {

		toutrnament := s.Find(".fpml-tournament a").Text()
		if toutrnament == "" {
			return
		}
		if rankInt >= 11 {
			return
		}

		left := s.Find("td .vs-align-left .markup-object span").Text()
		right := s.Find("td .vs-align-right .markup-object span").Text()
		toutrnamenttime := s.Find("td .countdowndate").Text()
		tTime := strings.Split(toutrnamenttime, " ")
		sTime := tTime[2] + "-" + month[tTime[1]] + "-" + tTime[0] + " " + tTime[3]
		timeObj1, _ := time.Parse("2006-1-2 15:04:05", sTime)
		timeObj1 = timeObj1.Local()

		fTime := timeObj1.Format("2006-01-02 15:04:05")

		str += "【" + strconv.Itoa(rankInt) + "】\n"
		str += "赛事：" + toutrnament + "\n"
		str += "比赛时间：" + fTime + "\n"
		str += "比赛队伍：" + left + " VS " + right
		str += "\n\n"

		rankInt++
	})

	msg.ReplyText(str)
}
