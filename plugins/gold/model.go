package gold

type ResponseGoldApi struct {
	Reason     string             `json:"reason"`
	ResultCode string             `json:"resultcode"`
	ErrorCode  int                `json:"error_code"`
	Result     []map[string]*List `json:"result"`
}

type List struct {
	Variety   string `json:"variety"`   // 品种
	Latestpri string `json:"latestpri"` // 最新价
	Openpri   string `json:"openpri"`   // 开盘价
	Maxpri    string `json:"maxpri"`    // 最高价
	Minpri    string `json:"minpri"`    // 最低价
	Limit     string `json:"limit"`     // 涨跌幅
	Yespri    string `json:"yespri"`    // 昨收价
	Totalvol  string `json:"totalvol"`  // 总成交量
	Time      string `json:"time"`      // 更新时间
}
