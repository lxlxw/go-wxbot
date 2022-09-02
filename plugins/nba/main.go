package nba

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
	"github.com/yqchilde/pkgs/log"
)

type NBA struct {
	engine.PluginMagic
	Enable bool   `yaml:"enable"`
	Url    string `yaml:"url"`
	Key    string `yaml:"key"`
}

var (
	keywords   = []string{"NBA赛程", "nba赛程", "nba联赛"}
	pluginInfo = &NBA{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 NBA赛程 => 获取NBA近期赛程",
			Commands: keywords,
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *NBA) OnRegister() {
}

func (p *NBA) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() {
			for _, v := range keywords {
				if msg.Content == v {
					getNBA(msg)
					return
				}
			}
		}
	}
}

func getNBA(msg *robot.Message) {

	var nbaConf NBA
	plugin.RawConfig.Unmarshal(&nbaConf)

	apiUrl := fmt.Sprintf("%s?key=%s", nbaConf.Url, nbaConf.Key)

	res, err := http.Get(apiUrl)
	if err != nil {
		log.Errorf("getNBA http get error: %v", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getNBA read body error: %v", err)
		return
	}

	var resp NBAApiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getNBA unmarshal error: %v", err)
		return
	}
	if resp.Error_code != 0 {
		log.Errorf("getNBA api error: %v", resp.Reason)
		return
	}
	var str = ""
	str += "【" + resp.Result.Duration + resp.Result.Title + "】\n\n"

	sTime := time.Now().Format("2006-01-02")

	nTime := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	var sList, nList string
	for _, v := range resp.Result.Matchs {
		if v.Date == sTime {
			for _, vv := range v.List {

				sList += "比赛时间：" + v.Date + " " + vv.Time_start + "\n"
				sList += "比赛队伍：" + vv.Team1 + " VS " + vv.Team2 + "\n"
				sList += "比赛状态：" + vv.Status_text
				if vv.Status == "3" || vv.Status == "2" {
					sList += "\n" + "比分：" + vv.Team1_score + ":" + vv.Team2_score
				}
				sList += "\n\n"
			}
		}
		if v.Date == nTime {
			for _, vv := range v.List {

				nList += "比赛时间：" + v.Date + " " + vv.Time_start + "\n"
				nList += "比赛队伍：" + vv.Team1 + " VS " + vv.Team2 + "\n"
				nList += "比赛状态：" + vv.Status_text
				if vv.Status == "3" || vv.Status == "2" {
					nList += "\n" + "比分：" + vv.Team1_score + ":" + vv.Team2_score
				}
				nList += "\n\n"
			}
		}
	}
	if len(sList) == 0 && len(nList) == 0 {
		str = "今明两天无NBA比赛"
	} else {
		str += sList
		str += nList
	}

	msg.ReplyText(str)
}
