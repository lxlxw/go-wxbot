package exchangerate

type ExApiResponse struct {
	Code int       `json:"error_code"`
	Msg  string    `json:"reason"`
	Data []*Result `json:"result"`
}

type Result struct {
	CurrencyF      string `json:"currencyF"`
	CurrencyF_Name string `json:"currencyF_Name"`
	CurrencyT      string `json:"currencyT"`
	CurrencyT_Name string `json:"currencyT_Name"`
	Exchange       string `json:"exchange"`
	UpdateTime     string `json:"updateTime"`
}
