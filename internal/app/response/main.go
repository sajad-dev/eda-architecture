package response

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Status  string      `json:"status"`
}
type ErrorResponse struct {
	Messages interface{} `json:"message"`
	Code     int         `json:"code"`
	Status   bool        `json:"status"`
}
