package food

type FoodApiResponse struct {
	Msg    string `json:"msg"`
	Code   int    `json:"code"`
	Result struct {
		List []*FoodList `json:"list"`
	} `json:"data"`
}

type FoodList struct {
	FoodId string      `json:"foodId"`
	Name   interface{} `json:"name"`
}

type FoodDetailResponse struct {
	Msg    string `json:"msg"`
	Code   int    `json:"code"`
	Result struct {
		Name             string `json:"name"`
		Calory           string `json:"calory"` // 热量
		CaloryUnit       string `json:"caloryUnit"`
		Protein          string `json:"protein"` // 蛋白质
		ProteinUnit      string `json:"proteinUnit"`
		Fat              string `json:"fat"` // 脂肪
		FatUnit          string `json:"fatUnit"`
		Carbohydrate     string `json:"carbohydrate"` // 碳水
		CarbohydrateUnit string `json:"carbohydrateUnit"`
		HealthLight      int    `json:"healthLight"`   //  健康等级 1 2 3 分别是推荐 适量 少吃
		HealthTips       string `json:"healthTips"`    // 健康描述
		HealthSuggest    string `json:"healthSuggest"` //	健康建议
	} `json:"data"`
}
