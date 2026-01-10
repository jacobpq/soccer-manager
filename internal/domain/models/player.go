package models

type Player struct {
	ID             int     `json:"id"`
	TeamID         int     `json:"team_id"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	Country        string  `json:"country"`
	Age            int     `json:"age"`
	Position       string  `json:"position"`
	Value          float64 `json:"value"`
	MarketPrice    float64 `json:"market_price,omitempty"`
	OnTransferList bool    `json:"on_transfer_list"`
}
