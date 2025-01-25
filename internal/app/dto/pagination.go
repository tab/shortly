package dto

// PaginatedResponse is a response with pagination
type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Page  int         `json:"page"`
	Per   int         `json:"per"`
	Total int         `json:"total"`
}
