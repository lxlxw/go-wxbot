package constellation

type ConTodayApiResponse struct {
	Code     int    `json:"error_code"`
	Name     string `json:"name"`
	QFriend  string `json:"QFriend"`
	Datetime string `json:"datetime"`
	Health   string `json:"health"`
	Love     string `json:"love"`
	Work     string `json:"work"`
	Money    string `json:"money"`
	Number   int    `json:"number"`
	Summary  string `json:"summary"`
	All      string `json:"all"`
	Color    string `json:"color"`
}

type CoWeekApiResponse struct {
	Code   int    `json:"error_code"`
	Name   string `json:"name"`
	Date   string `json:"date"`
	Health string `json:"health"`
	Love   string `json:"love"`
	Money  string `json:"money"`
	Weekth int    `json:"weekth"`
	Work   string `json:"work"`
	All    string `json:"all"`
}

type CoMonthApiResponse struct {
	Code   int    `json:"error_code"`
	Name   string `json:"name"`
	Date   string `json:"date"`
	Health string `json:"health"`
	Love   string `json:"love"`
	Money  string `json:"money"`
	Weekth int    `json:"weekth"`
	Work   string `json:"work"`
}

type CoYearApiResponse struct {
	Code int    `json:"error_code"`
	Name string `json:"name"`
	Date string `json:"date"`
	Mima struct {
		Info string   `json:"info"`
		Text []string `json:"text"`
	} `json:"mima"`
	Career      []string `json:"career"`
	Love        []string `json:"love"`
	Health      []string `json:"health"`
	Finance     []string `json:"finance"`
	LuckeyStone string   `json:"luckeyStone"`
}
