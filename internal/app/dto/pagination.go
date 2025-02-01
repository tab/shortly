package dto

// PaginatedResponse is a response with pagination
type PaginatedResponse struct {
	// Data is the data to be returned
	Data interface{} `json:"data"`
	// Page is the current page
	Page int `json:"page"`
	// Per is the number of items per page
	Per int `json:"per"`
	// Total is the total number of items
	Total int `json:"total"`
}
