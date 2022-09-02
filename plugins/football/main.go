package football

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

type Football struct {
	engine.PluginMagic
	Enable bool   `yaml:"enable"`
	Url    string `yaml:"url"`
	Key    string `yaml:"key"`
}

var footMap = map[string]string{
	"英超赛程": "yingchao",
	"西甲赛程": "xijia",
	"德甲赛程": "dejia",
	"意甲赛程": "yijia",
	"法甲赛程": "fajia",
	"中超赛程": "zhongchao",
	"英超联赛": "yingchao",
	"西甲联赛": "xijia",
	"德甲联赛": "dejia",
	"意甲联赛": "yijia",
	"法甲联赛": "fajia",
	"中超联赛": "zhongchao",
}

var (
	keywords   = []string{"英超赛程", "西甲赛程", "德甲赛程", "意甲赛程", "法甲赛程", "中超赛程", "英超联赛", "西甲联赛", "德甲联赛", "意甲联赛", "法甲联赛", "中超联赛"}
	pluginInfo = &Football{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {name}赛程 => 获取五大联赛当天赛程 | 示例：英超赛程",
			Commands: keywords,
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Football) OnRegister() {
}

func (p *Football) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() {
			for _, v := range keywords {
				if msg.Content == v {
					getFootball(msg)
					return
				}
			}
		}
	}
}

func getFootball(msg *robot.Message) {

	var footConf Football
	plugin.RawConfig.Unmarshal(&footConf)

	apiUrl := fmt.Sprintf("%s?type=%s&key=%s", footConf.Url, footMap[msg.Content], footConf.Key)

	res, err := http.Get(apiUrl)
	if err != nil {
		log.Errorf("getFootball http get error: %v", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getFootball read body error: %v", err)
		return
	}

	var resp FootApiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getFootball unmarshal error: %v", err)
		return
	}
	if resp.Error_code != 0 {
		log.Errorf("getFootball api error: %v", resp.Reason)
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
		str = "今日无该联赛比赛"
	} else {
		str += sList
		str += nList
	}

	msg.ReplyText(str)
}
