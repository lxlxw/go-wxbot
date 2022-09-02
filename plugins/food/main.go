package food

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
	"github.com/yqchilde/pkgs/log"
)

type Food struct {
	engine.PluginMagic
	Enable    bool   `yaml:"enable"`
	Url1      string `yaml:"url1"`
	Url2      string `yaml:"url2"`
	AppId     string `yaml:"appId"`
	AppSecret string `yaml:"appSecret"`
}

var (
	keywords   = []string{"营养成分", "热量", "能量", "脂肪", "蛋白质"}
	pluginInfo = &Food{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {food}营养成分 => 获取食物营养成分 || 示例：香蕉营养成分",
			Commands: keywords,
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Food) OnRegister() {
}

func (p *Food) OnEvent(msg *robot.Message) {
	if msg != nil {
		for _, v := range keywords {
			if strings.Contains(msg.Content, v) {
				getFood(msg, v)
				return
			}
		}
	}
}

func getFood(msg *robot.Message, keyword string) {

	var foodConf Food
	plugin.RawConfig.Unmarshal(&foodConf)

	foodName := strings.Trim(msg.Content, keyword)

	apiUrl1 := fmt.Sprintf("%s?keyword=%s&page=1&app_id=%s&app_secret=%s", foodConf.Url1, foodName, foodConf.AppId, foodConf.AppSecret)

	res, err := http.Get(apiUrl1)
	if err != nil {
		log.Errorf("getFood http get error: %v", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getFood read body error: %v", err)
		return
	}

	var resp FoodApiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getCook unmarshal error: %v", err)
		return
	}
	if resp.Code != 1 {
		log.Errorf("getFood api error: %v", resp.Msg)
		return
	}

	var str string

	if len(resp.Result.List) <= 0 {
		log.Errorf("getFood api error: %v", resp.Msg)
		str = "未找到该食物营养成分表"
		msg.ReplyText(str)
		return
	}
	var foodId string
	for _, v := range resp.Result.List {
		if v.Name == foodName {
			foodId = v.FoodId
			continue
		}
	}

	// get food detail
	apiUrl2 := fmt.Sprintf("%s?foodId=%s&page=1&app_id=%s&app_secret=%s", foodConf.Url2, foodId, foodConf.AppId, foodConf.AppSecret)

	res2, err := http.Get(apiUrl2)
	if err != nil {
		log.Errorf("getFood http get error: %v", err)
		return
	}
	defer res2.Body.Close()

	body2, err := io.ReadAll(res2.Body)
	if err != nil {
		log.Errorf("getFood read body error: %v", err)
		return
	}

	var resp2 FoodDetailResponse
	if err := json.Unmarshal(body2, &resp); err != nil {
		log.Errorf("getCook unmarshal error: %v", err)
		return
	}
	if resp2.Code != 1 {
		log.Errorf("getFood api error: %v", resp.Msg)
		return
	}

	str += "【" + resp2.Result.Name + "】" + "营养成分表" + "\n\n"

	str += "热量：" + resp2.Result.Calory + resp2.Result.CaloryUnit + "\n"
	str += "蛋白质：" + resp2.Result.Protein + resp2.Result.ProteinUnit + "\n"
	str += "碳水化合物：" + resp2.Result.Carbohydrate + resp2.Result.CarbohydrateUnit + "\n"
	str += "脂肪：" + resp2.Result.Fat + resp2.Result.FatUnit + "\n"
	str += "健康描述：" + resp2.Result.HealthTips + "\n"
	str += "健康建议：" + resp2.Result.HealthSuggest + "\n"

	msg.ReplyText(str)

}
