package models

type Team struct {
	ID      int     `json:"id"`
	UserID  int     `json:"user_id"`
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Budget  float64 `json:"budget"`
	Value   float64 `json:"total_value"`
}
