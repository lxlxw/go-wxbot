package food

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	notFound   = "未找到该食物营养成分表"
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
			if strings.Contains(msg.Content, v) && msg.Content != notFound {
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

	data := url.Values{}
	data.Add("keyword", foodName)
	data.Add("page", "1")
	data.Add("app_id", foodConf.AppId)
	data.Add("app_secret", foodConf.AppSecret)

	u, _ := url.ParseRequestURI("https://www.mxnzp.com")

	u.Path = "/api/food_heat/food/search/"
	u.RawQuery = data.Encode()
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return
	}
	res, _ := client.Do(req)

	defer res.Body.Close()

	body1, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getFood read body error: %v", err)
		return
	}

	var resp FoodApiResponse
	if err := json.Unmarshal(body1, &resp); err != nil {
		log.Errorf("getFood read body error: %v", err)
		return
	}
	if resp.Code != 1 {
		log.Errorf("getFood api error: %v", resp.Msg)
		return
	}

	var str string

	if len(resp.Result.List) <= 0 {
		log.Errorf("getFood api error: %v", resp.Msg)
		str = notFound
		msg.ReplyText(str)
		return
	}
	var foodId string
	for _, v := range resp.Result.List {
		if v.Name == foodName {
			foodId = v.FoodId
			break
		}
	}

	data2 := url.Values{}
	data2.Add("foodId", foodId)
	data2.Add("page", "1")
	data2.Add("app_id", foodConf.AppId)
	data2.Add("app_secret", foodConf.AppSecret)

	u2, _ := url.ParseRequestURI("https://www.mxnzp.com")

	u2.Path = "/api/food_heat/food/details/"
	u2.RawQuery = data2.Encode()
	urlStr2 := fmt.Sprintf("%v", u2)
	client2 := &http.Client{}
	req2, err := http.NewRequest("GET", urlStr2, nil)
	if err != nil {
		return
	}
	res2, _ := client2.Do(req2)
	defer res2.Body.Close()

	body2, err := io.ReadAll(res2.Body)
	if err != nil {
		log.Errorf("getFood read body error: %v", err)
		return
	}

	var resp2 FoodDetailResponse
	if err := json.Unmarshal(body2, &resp2); err != nil {
		log.Errorf("getCook unmarshal error: %v", err)
		return
	}
	if resp2.Code != 1 {
		log.Errorf("getFood api error: %v", resp2.Msg)
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
