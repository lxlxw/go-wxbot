package weibo

import (
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/yqchilde/pkgs/log"

	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
)

type Weibo struct {
	engine.PluginMagic
	Enable bool   `yaml:"enable"`
	Url    string `yaml:"url"`
	Key    string `yaml:"key"`
}

const HOSTNAME = "https://s.weibo.com"

var (
	keyword    = "微博热搜"
	pluginInfo = &Weibo{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 微博热搜 => 获取微博top10",
			Commands: []string{keyword},
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Weibo) OnRegister() {
}

func (p *Weibo) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() && msg.Content == keyword {
			getWeiboDetail(msg)
		}
	}
}

func getWeiboDetail(msg *robot.Message) {

	// TODO 获取微博top10
	content, err := getWeiboTop()
	if err != nil {
		log.Errorf("getWeiboDetail http get error: %v", err)
		return
	}

	msg.ReplyText(content)
}

func getWeiboTop() (string, error) {

	var weiboConf Weibo
	plugin.RawConfig.Unmarshal(&weiboConf)

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://s.weibo.com/top/summary?cate=realtimehot", nil)

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36")
	req.Header.Add("Cookie", "SUB=_2AkMVWDYUf8NxqwJRmP0Sz2_hZYt2zw_EieKjBMfPJRMxHRl-yj9jqkBStRB6PtgY-38i0AF7nDAv8HdY1ZwT3Rv8B5e5; SUBP=0033WrSXqPxfM72-Ws9jqgMF55529P9D9WFencmWZyNhNlrzI6f0SiqP")
	if err != nil {
		return "", err
	}
	res, _ := client.Do(req)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	var str = ""
	var rankInt = 0
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		if rankInt >= 10 {
			return
		}
		redu := s.Find(".td-03 i").Text()
		if redu != "商" && redu != "首发" && redu != "上新" {

			href, _ := s.Find(".td-02 a").Attr("href")
			herfText := s.Find(".td-02 a").Text()
			redu := s.Find(".td-03 i").Text()
			rank := s.Find(".ranktop").Text()
			if rank == "" {
				rank = redu
			}

			w_url := HOSTNAME + href

			apiUrl := fmt.Sprintf("%s?url=%s&key=%s", weiboConf.Url, w_url, weiboConf.Key)

			res, err := http.Get(apiUrl)
			if err != nil {
				log.Errorf("getWeiboDetail http get error: %v", err)
				return
			}
			body, err := io.ReadAll(res.Body)
			if err != nil {
				log.Errorf("getWeiboDetail read body error: %v", err)
				return
			}
			str += fmt.Sprintf(`%s、%s：%s`, rank, herfText, string(body))
			str += "\n"
			rankInt++
		}
	})
	return str, nil
}
