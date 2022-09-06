package weather

type ResponseWeatherLocaltion struct {
	Code     string         `json:"code"` // 状态码
	Location []LocationList `json:"location"`
}
type LocationList struct {
	Name      string `json:"name"`      // 地区/城市名称
	ID        string `json:"id"`        // 地区/城市ID
	Lat       string `json:"lat"`       // 地区/城市纬度
	Lon       string `json:"lon"`       // 地区/城市经度
	Adm2      string `json:"adm2"`      // 地区/城市的上级行政区划名称
	Adm1      string `json:"adm1"`      // 地区/城市所属一级行政区域
	Country   string `json:"country"`   // 地区/城市所属国家名称
	Tz        string `json:"tz"`        // 地区/城市所在时区
	Utcoffset string `json:"utcOffset"` // 地区/城市目前与UTC时间偏移的小时数
	Isdst     string `json:"isDst"`     // 地区/城市是否当前处于夏令时 1 表示当前处于夏令时 0 表示当前不是夏令时
	Type      string `json:"type"`      // 地区/城市的属性
	Rank      string `json:"rank"`      // 地区评分
	Fxlink    string `json:"fxLink"`    // 该地区的天气预报网页链接，便于嵌入你的网站或应用
}

type ResponseWeatherActual struct {
	Code       string `json:"code"`       // 状态码
	Updatetime string `json:"updateTime"` // 当前API的最近更新时间
	Fxlink     string `json:"fxLink"`     // 当前数据的响应式页面，便于嵌入网站或应用
	Now        Now    `json:"now"`
}
type Now struct {
	Obstime   string `json:"obsTime"`   // 数据观测时间
	Temp      string `json:"temp"`      // 温度，默认单位：摄氏度
	Feelslike string `json:"feelsLike"` // 体感温度，默认单位：摄氏度
	Icon      string `json:"icon"`      // 天气状况和图标的代码，图标可通过天气状况和图标下载
	Text      string `json:"text"`      // 天气状况的文字描述，包括阴晴雨雪等天气状态的描述
	Wind360   string `json:"wind360"`   // 风向360角度
	Winddir   string `json:"windDir"`   // 风向
	Windscale string `json:"windScale"` // 风力等级
	Windspeed string `json:"windSpeed"` // 风速，公里/小时
	Humidity  string `json:"humidity"`  // 相对湿度，百分比数值
	Precip    string `json:"precip"`    // 当前小时累计降水量，默认单位：毫米
	Pressure  string `json:"pressure"`  // 大气压强，默认单位：百帕
	Vis       string `json:"vis"`       // 能见度，默认单位：公里
	Cloud     string `json:"cloud"`     // 云量，百分比数值。可能为空
	Dew       string `json:"dew"`       // 露点温度。可能为空
}

type ResponseWeatherLife struct {
	Code       string  `json:"code"`       // 状态码
	Updatetime string  `json:"updateTime"` // 当前API的最近更新时间
	Fxlink     string  `json:"fxLink"`     // 当前数据的响应式页面，便于嵌入网站或应用
	Daily      []Daily `json:"daily"`
}

type Daily struct {
	Date     string `json:"date"`     //     预报日期
	Type     string `json:"type"`     // 生活指数类型ID
	Name     string `json:"name"`     // 生活指数类型的名称
	Level    string `json:"level"`    // 生活指数预报等级
	Category string `json:"category"` // 生活指数预报级别名称
	Text     string `json:"text"`     // 生活指数预报的详细描述，可能为空
}

type ResponseWeatherAir struct {
	Code       string `json:"code"`       // 状态码
	Updatetime string `json:"updateTime"` // 当前API的最近更新时间
	Fxlink     string `json:"fxLink"`     // 当前数据的响应式页面，便于嵌入网站或应用
	Now        AirNow `json:"now"`
}
type AirNow struct {
	Pubtime  string `json:"pubTime"`  // 空气质量数据发布时间
	Aqi      string `json:"aqi"`      // 空气质量指数
	Level    string `json:"level"`    // 空气质量指数等级
	Category string `json:"category"` // 空气质量指数级别
	Primary  string `json:"primary"`  // 空气质量的主要污染物，空气质量为优时，返回值为NA
	Pm10     string `json:"pm10"`     // PM10
	Pm2P5    string `json:"pm2p5"`    // PM2.5
	No2      string `json:"no2"`      // 二氧化氮
	So2      string `json:"so2"`      // 二氧化硫
	Co       string `json:"co"`       // 一氧化碳
	O3       string `json:"o3"`       // 臭氧
}

type ResponseWeatherWarning struct {
	Code       string    `json:"code"`       // 状态码
	Updatetime string    `json:"updateTime"` // 当前API的最近更新时间
	Fxlink     string    `json:"fxLink"`     // 当前数据的响应式页面，便于嵌入网站或应用
	Warning    []Warning `json:"warning"`
}
type Warning struct {
	ID            string `json:"id"`            // 预警发布单位，可能为空
	Sender        string `json:"sender"`        // 预警发布时间
	Pubtime       string `json:"pubTime"`       // 预警发布时间
	Title         string `json:"title"`         // 预警信息标题
	Starttime     string `json:"startTime"`     // 预警开始时间，可能为空
	Endtime       string `json:"endTime"`       // 预警结束时间，可能为空
	Status        string `json:"status"`        // 预警信息的发布状态
	Severity      string `json:"severity"`      // 预警严重等级 Unknown、Minor、Moderate、Severe和Extreme
	Severitycolor string `json:"severityColor"` // 预警严重等级颜色，可能为空
	Type          string `json:"type"`          // 预警类型ID
	Typename      string `json:"typeName"`      // 预警类型名称
	Text          string `json:"text"`          // 预警详细文字描述
}

type ResponseWeatherSun struct {
	Code       string `json:"code"`       // 状态码
	Updatetime string `json:"updateTime"` // 当前API的最近更新时间
	Fxlink     string `json:"fxLink"`     // 当前数据的响应式页面，便于嵌入网站或应用
	Sunrise    string `json:"sunrise"`    // 日出时间
	Sunset     string `json:"sunset"`     // 日落时间
}
