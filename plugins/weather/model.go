package weather

type WeatherApiResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Address       string `json:"address"`
		CityCode      string `json:"cityCode"`
		Temp          string `json:"temp"`
		Weather       string `json:"weather"`
		WindDirection string `json:"windDirection"`
		WindPower     string `json:"windPower"`
		Humidity      string `json:"humidity"`
		ReportTime    string `json:"reportTime"`
	} `json:"data"`
}
