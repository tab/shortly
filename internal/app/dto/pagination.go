package dto

type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Page  int         `json:"page"`
	Per   int         `json:"per"`
	Total int         `json:"total"`
}
