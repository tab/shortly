package dto

// StatsResponse is a response for stats request
type StatsResponse struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}
