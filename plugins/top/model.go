package top

type TopApiResponse struct {
	Msg    string `json:"msg"`
	Code   int    `json:"code"`
	Result struct {
		Name        string  `json:"name"`
		Last_update string  `json:"last_update"`
		List        []*List `json:"list"`
	} `json:"data"`
}

type List struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	Other string `json:"other"`
}
