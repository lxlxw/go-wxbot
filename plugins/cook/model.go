package cook

type CookApiResponse struct {
	Msg    string `json:"msg"`
	Code   string `json:"code"`
	Result struct {
		Status interface{} `json:"status"`
		Msg    string      `json:"msg"`
		Result *CookResult `json:"result"`
	} `json:"result"`
}

type CookResult struct {
	Num  interface{} `json:"num"`
	List []*CookList `json:"list"`
}

type CookList struct {
	Id          interface{} `json:"id"`
	Classid     interface{} `json:"classid"`
	Name        string      `json:"name"`
	Peoplenum   string      `json:"peoplenum"`
	Preparetime string      `json:"preparetime"`
	Cookingtime string      `json:"cookingtime"`
	Tag         string      `json:"tag"`
	Material    []*Material `json:"material"`
	Process     []*Process  `json:"process"`
}

type Material struct {
	Mname    string `json:"mname"`
	Nametype string `json:"type"`
	Amount   string `json:"amount"`
}

type Process struct {
	Pcontent string `json:"pcontent"`
}
