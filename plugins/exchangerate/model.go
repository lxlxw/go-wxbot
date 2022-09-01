package exchangerate

type ExApiResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Name     string `json:"name"`
		NameDesc string `json:"nameDesc"`
		From     string `json:"from"`
		To       string `json:"to"`
		Price    string `json:"price"`
	} `json:"data"`
}
