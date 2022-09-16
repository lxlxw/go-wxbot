package jx3

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/FloatTech/AnimeAPI/pixiv"
	"github.com/FloatTech/floatbox/web"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/lxlxw/go-wxbot/engine"
	"github.com/lxlxw/go-wxbot/engine/robot"
)

type Jx3 struct {
	engine.PluginMagic
	Enable bool `yaml:"enable"`
}

var (
	jc   = "剑纯宏"
	qc   = "气纯宏"
	tc   = "天策宏"
	tc_o = "傲血宏"
	tlh  = "铁牢宏"
	yj   = "易经宏"
	xs   = "洗髓宏"
	dj   = "毒经宏"
	jy   = "惊羽宏"
	jy_o = "鲸鱼宏"
	tl   = "天罗宏"
	tl_o = "田螺宏"
	ws   = "问水宏"
	ws_o = "藏剑宏"
	bx   = "冰心宏"
	hj   = "花间宏"
)

var (
	keywords   = []string{jc, qc, tc, yj, xs, dj, jy, jy_o, tl, tl_o, ws, ws_o, bx, hj}
	pluginInfo = &Jx3{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {门派}宏 => 获取jx3缘起门派宏 || 示例：剑纯宏",
			Commands: keywords,
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (p *Jx3) OnRegister() {
}

func (p *Jx3) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.IsText() {
			for _, v := range keywords {
				if msg.Content == v {
					getJx3Detail(msg, msg.Content)
					return
				}
			}
			if msg.Content == "/roll" {
				ranInt := rand.Intn(100)
				msg.ReplyText(strconv.Itoa(ranInt))
				return
			}
			if msg.Content == "吃什么" {
				getEat(msg)
				return
			}

			if strings.Contains(msg.Content, "来一张") {
				getImages(msg)
				return
			}
		}
	}
}

type resultjson struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    struct {
		Illusts []struct {
			ID          int64  `json:"id"`
			Title       string `json:"title"`
			AltTitle    string `json:"altTitle"`
			Description string `json:"description"`
			Type        int64  `json:"type"`
			CreateDate  string `json:"createDate"`
			UploadDate  string `json:"uploadDate"`
			Sanity      int64  `json:"sanity"`
			Width       int64  `json:"width"`
			Height      int64  `json:"height"`
			PageCount   int64  `json:"pageCount"`
			Tags        []struct {
				Name        string `json:"name"`
				Translation string `json:"translation"`
			} `json:"tags"`
			Statistic struct {
				Bookmarks int64 `json:"bookmarks"`
				Likes     int64 `json:"likes"`
				Comments  int64 `json:"comments"`
				Views     int64 `json:"views"`
			} `json:"statistic"`
			Image string `json:"image"`
		} `json:"illusts"`
		Scores  []float64 `json:"scores"`
		HasNext bool      `json:"has_next"`
	} `json:"data"`
}

// soutuapi 请求api
func soutuapi(keyword string) (r resultjson, err error) {
	var data []byte
	data, err = web.RequestDataWith(&http.Client{Transport: &http3.RoundTripper{}},
		"https://api.pixivel.moe/v2/pixiv/illust/search/"+url.QueryEscape(keyword)+"?page=0",
		"GET",
		"https://pixivel.moe/",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36",
	)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &r)
	if err == nil && r.Error {
		err = errors.New(r.Message)
	}
	return
}

func getImages(msg *robot.Message) {
	keyword := strings.Trim(msg.Content, "来一张")
	soutujson, err := soutuapi(keyword)
	if err != nil {
		return
	}
	rannum := rand.Intn(len(soutujson.Data.Illusts))
	il := soutujson.Data.Illusts[rannum]
	illust, err := pixiv.Works(il.ID)
	if err != nil {
		return
	}
	f := illust.Path(0)
	if illust.DownloadToCache(0) != nil {
		fmt.Println(err)
		return
	}

	file, err := os.Open(f)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	msg.ReplyFile(file)

	return
}

type EatApiResponse struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	List []string `json:"data"`
}

func getEat(msg *robot.Message) {

	apiUrl := fmt.Sprintf("%s?size=7", "https://eolink.o.apispace.com/eat222/api/v1/forward/chishenme")
	println(apiUrl)
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)

	req.Header.Add("X-APISpace-Token", "lfl0e6t5db07d55ecwdvr1dys7dixumg")
	req.Header.Add("Authorization-Type", "apikey")
	if err != nil {
		return
	}
	res, _ := client.Do(req)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	var resp EatApiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Println("err", err)
		return
	}
	if resp.Code != 200 {
		return
	}
	var str string
	if len(resp.List) <= 0 {
		return
	}

	str += "今天吃什么，我来给你建议：" + "\n"
	for k, v := range resp.List {
		str += strconv.Itoa(k+1) + " - " + v + "\n"
	}
	msg.ReplyText(str)

}

func getJx3Detail(msg *robot.Message, keyword string) {

	detail := getJx3Reply(keyword)
	msg.ReplyText(detail)
}

func getJx3Reply(keyword string) string {
	var detail string
	if keyword == jc {
		detail = `/cast 猛虎下山` + "\n" +
			`/cast 凭虚御风` + "\n" +
			`/cast [qidian<6] 剑飞惊天` + "\n" +
			`/fcast [qidian>7|tbufftime:叠刃<7] 无我无剑` + "\n" +
			`/fcast 天地无极` + "\n" +
			`/fcast 三环套月` + "\n" +
			`/cast 凝神聚气`
	} else if keyword == qc {
		detail = `/fcast [nobuff:破苍穹·期声] 破苍穹` + "\n" +
			`/cast 凭虚御风` + "\n" +
			`/fcast [qidian>7|tbuff:无形&qidian>6] 两仪化形` + "\n" +
			`/fcast 四象轮回` + "\n" +
			`/cast 凝神聚气`
	} else if keyword == bx {
		detail = `/cast [tnobuff:急曲] 剑主天地` + "\n" +
			`/cast [nobuff:剑神无我] 剑神无我` + "\n" +
			`/cast [tbuff:急曲>=2] 剑气长江` + "\n" +
			`/cast [bufftime:碎冰>=2&tbuff:急曲>=2] 江海凝光` + "\n" +
			`/cast 玳弦急曲`
	} else if keyword == tc || keyword == tc_o {
		detail = `/cast 猛虎下山` + "\n" +
			`/cast 啸如虎` + "\n" +
			`/cast [tnobuff:致残] 龙吟` + "\n" +
			`/cast [tnobuff:流血] 破风` + "\n" +
			`/cast [tbufftime:流血<2|tbuff:致残] 灭` + "\n" +
			`/cast 龙牙` + "\n" +
			`/cast 霹雳` + "\n" +
			`/cast 穿云`
	} else if keyword == yj {
		detail = `/cast 佛心诀` + "\n" +
			`/cast 猛虎下山` + "\n" +
			`/cast [qidian>2] 金刚怒目` + "\n" +
			`/cast [qidian>2] 拿云式` + "\n" +
			`/cast [qidian>2] 韦陀献杵` + "\n" +
			`/cast 守缺式` + "\n" +
			`/cast 横扫六合` + "\n" +
			`/cast 捣虚式` + "\n" +
			`/cast 普渡四方`
	} else if keyword == xs {
		detail = `/cast [qidian>2] 袖纳乾坤` + "\n" +
			`/cast [tbuff:立地成佛>2&qidian>2] 灵山施雨` + "\n" +
			`/cast [qidian>2] 立地成佛` + "\n" +
			`/cast 大狮子吼` + "\n" +
			`/cast 横扫六合` + "\n" +
			`/cast 普渡四方` + "\n" +
			`/cast 擒龙诀`
	} else if keyword == dj {
		detail = `/cast 蛊虫献祭` + "\n" +
			`/cast 夺命蛊` + "\n" +
			`/cast 灵蛇引` + "\n" +
			`/cast 圣蝎引` + "\n" +
			`/cast 攻击` + "\n" +
			`/cast 幻击` + "\n" +
			`/cast 百足` + "\n" +
			`/cast [tnobuff:蛇影] 蛇影` + "\n" +
			`/cast 蟾啸` + "\n" +
			`/cast 蝎心` + "\n" + "\n" +
			`进战斗前先召蝎子，最好战斗前三十秒以上召唤出来，手动狂暴！`

	} else if keyword == jy || keyword == jy_o {
		detail = `/fcast [tbuff:千疮百孔] 暴雨梨花针` + "\n" +
			`/cast [energy<30] 连环弩` + "\n" +
			`/cast 孔雀翎` + "\n" +
			`/cast 猛虎下山` + "\n" +
			`/cast 逐星箭` + "\n" +
			`/cast [nobuff:追命无声] 夺魄箭` + "\n" +
			`/cast 追命箭`
	} else if keyword == tl || keyword == tl_o {
		detail = `/cast 天绝地灭` + "\n" +
			`/cast [tnobuff:化血镖] 化血镖` + "\n" +
			`/cast [energy>80] 暗藏杀机` + "\n" +
			`/cast [nobuff:心无旁骛] 暴雨梨花针` + "\n" +
			`/cast [nobuff:奥妙&nobuff:心无旁骛] 孔雀翎` + "\n" +
			`/cast 蚀肌弹`
	} else if keyword == ws || keyword == ws_o {
		detail = `/cast [rage<40] 潮鸣弦` + "\n" +
			`/cast [rage<20] 莺鸣柳` + "\n" +
			`/cast [rage<20] 雪断桥` + "\n" +
			`/cast [rage<20] 云栖松` + "\n" +
			`/cast 猛虎下山` + "\n" +
			`/cast 断潮` + "\n" +
			`/fcast 云飞玉皇` + "\n" +
			`/cast 夕照雷峰`
	} else if keyword == hj {
		detail = `/cast [tnobuff:兰摧玉折] 兰摧玉折` + "\n" +
			`/cast [last_skill=兰摧玉折] 钟林毓秀` + "\n" +
			`/cast [tnobuff:商阳指] 商阳指` + "\n" +
			`/cast [buff:满雪=2] 快雪时晴` + "\n" +
			`/cast 玉石俱焚` + "\n" +
			`/cast 阳明指`
	} else if keyword == tlh {
		detail = `/cast 啸如虎` + "\n" +
			`/cast 灭` + "\n" +
			`/cast 破风`
	}
	return detail
}
