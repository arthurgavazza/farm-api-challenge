package models


type PaginatedResponse[T any] struct {
	Items       []T   `json:"items"`
	TotalCount  int   `json:"total_count"`
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
}
