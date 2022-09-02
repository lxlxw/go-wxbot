package football

type FootApiResponse struct {
	Reason     string `json:"reason"`
	Error_code int    `json:"error_code"`
	Result     struct {
		Title    string    `json:"title"`
		Duration string    `json:"duration"`
		Matchs   []*Matchs `json:"matchs"`
	} `json:"result"`
}

type Matchs struct {
	Date string       `json:"date"`
	Week string       `json:"week"`
	List []*MatchList `json:"list"`
}

type MatchList struct {
	Time_start  string `json:"time_start"`
	Status      string `json:"status"`
	Status_text string `json:"status_text"`
	Team1       string `json:"team1"`
	Team2       string `json:"team2"`
	Team1_score string `json:"team1_score"`
	Team2_score string `json:"team2_score"`
}
