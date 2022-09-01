package main

import (
	"context"

	"github.com/lxlxw/go-wxbot/engine"
	_ "github.com/lxlxw/go-wxbot/plugins/chat"         // 聊天机器人
	_ "github.com/lxlxw/go-wxbot/plugins/covid19"      // 疫情情况
	_ "github.com/lxlxw/go-wxbot/plugins/crazykfc"     // 肯德基疯狂星期四
	_ "github.com/lxlxw/go-wxbot/plugins/crypto"       // 加密货币
	_ "github.com/lxlxw/go-wxbot/plugins/exchangerate" // 实时汇率
	_ "github.com/lxlxw/go-wxbot/plugins/gold"         // 黄金
	_ "github.com/lxlxw/go-wxbot/plugins/jx3"          // jx3
	_ "github.com/lxlxw/go-wxbot/plugins/lottery"      // 彩票
	_ "github.com/lxlxw/go-wxbot/plugins/oil"          // 油价
	_ "github.com/lxlxw/go-wxbot/plugins/pinyinsuoxie" // 拼音缩写翻译
	_ "github.com/lxlxw/go-wxbot/plugins/stocks"       // 股票信息
	_ "github.com/lxlxw/go-wxbot/plugins/weather"      // 天气预报
	_ "github.com/lxlxw/go-wxbot/plugins/weibo"        // 微博热搜
	_ "github.com/lxlxw/go-wxbot/plugins/worldcup"     // 世界杯
)

func main() {
	ctx := context.Background()
	engine.Run(ctx, "config.yaml")
}
