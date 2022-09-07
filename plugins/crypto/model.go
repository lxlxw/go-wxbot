package crypto

type CryptoApiResponse struct {
	Code string        `json:"code"`
	Msg  string        `json:"msg"`
	List []*CryptoList `json:"data"`
}

type CryptoList struct {
	Symbol       string  `json:"symbol"`
	Price        float64 `json:"price"`
	ExchangeName string  `json:"exchangeName"`
}
