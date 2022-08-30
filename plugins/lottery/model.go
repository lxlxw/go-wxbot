package lottery

type LotteryApiResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Code     string `json:"code"`
		OpenCode string `json:"openCode"`
		Expect   string `json:"expect"`
		Name     string `json:"name"`
		Time     string `json:"time"`
	} `json:"data"`
}
