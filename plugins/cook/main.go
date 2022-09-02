package cook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
	"github.com/yqchilde/pkgs/log"
)

type Cook struct {
	engine.PluginMagic
	Enable bool   `yaml:"enable"`
	Url    string `yaml:"url"`
	Key    string `yaml:"key"`
}

var (
	keyword    = "菜谱"
	pluginInfo = &Cook{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {cook}菜谱 => 获取菜谱做法 || 示例：红烧排骨菜谱",
			Commands: []string{keyword},
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Cook) OnRegister() {
}

func (p *Cook) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() && strings.Contains(msg.Content, pluginInfo.Commands[0]) {
			getCook(msg)
		}
	}
}

func getCook(msg *robot.Message) {

	var cookConf Cook
	plugin.RawConfig.Unmarshal(&cookConf)

	cookName := strings.Trim(msg.Content, keyword)

	apiUrl := fmt.Sprintf("%s?keyword=%s&num=1&appkey=%s", cookConf.Url, cookName, cookConf.Key)

	res, err := http.Get(apiUrl)
	if err != nil {
		log.Errorf("getCook http get error: %v", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getCook read body error: %v", err)
		return
	}

	var resp CookApiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getCook unmarshal error: %v", err)
		return
	}
	if resp.Code != "10000" {
		log.Errorf("getCook api error: %v", resp.Msg)
		return
	}

	var str string

	if len(resp.Result.Result.List) <= 0 {
		log.Errorf("getCook api error: %v", resp.Msg)
		str = "未找到该菜谱做法"
		msg.ReplyText(str)
		return
	}

	List0 := resp.Result.Result.List[0]

	str += "【" + List0.Name + "】\n\n"
	str += "菜谱类型：" + List0.Tag + "\n"
	str += "预计烹饪时间：" + List0.Cookingtime + "\n"

	str += "所需材料：" + "\n"
	for _, v := range List0.Material {
		str += v.Mname + "：" + v.Amount + "\n"
	}

	str += "\n" + "烹饪步骤：" + "\n"
	for k, v := range List0.Process {
		str += "【" + strconv.Itoa(k+1) + "】" + v.Pcontent + "\n"
	}
	msg.ReplyText(str)

}
