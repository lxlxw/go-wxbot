package crypto

type CryptoApiResponse struct {
	Symbol           string  `json:"symbol"`
	Last_trade_price float64 `json:"last_trade_price"`
}
