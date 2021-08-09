package response

type Response struct {
	Code    StatusCode  `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
