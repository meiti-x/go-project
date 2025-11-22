package api

type APIResponse[TData any] struct {
	BaseResponse
	Data TData `json:"data"`
}

type BaseResponse struct {
	Status  string   `json:"status"`
	Message []string `json:"message"`
}

type ApiPaginateResponse[TData any] struct {
	TotalPage   uint    `json:"total_page"`
	CurrentPage uint    `json:"current_page"`
	Items       []TData `json:"items"`
}
