package main

import (
	"context"

	"github.com/lxlxw/go-wxbot/engine"
	_ "github.com/lxlxw/go-wxbot/plugins/crazykfc"     // 肯德基疯狂星期四骚话
	_ "github.com/lxlxw/go-wxbot/plugins/pinyinsuoxie" // 拼音缩写翻译
	_ "github.com/lxlxw/go-wxbot/plugins/weather"      // 天气预报
)

func main() {
	ctx := context.Background()
	engine.Run(ctx, "config.yaml")
}
