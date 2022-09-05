package oil

type ResponseOilApi struct {
	Reason    string  `json:"reason"`
	ErrorCode int     `json:"error_code"`
	Result    []*List `json:"result"`
}

type List struct {
	City   string `json:"city"`
	Oil92h string `json:"92h"`
	Oil95h string `json:"95h"`
	Oil98h string `json:"98h"`
	Oil0h  string `json:"0h"`
}
