package resp

type Response struct {
	*BaseResponse
	Data interface{} `json:"data,omitempty"`
}

type BaseResponse struct {
	Code    int    `json:"code" example:"10000"`
	Message string `json:"msg" example:"success"`
	Detail  string `json:"detail,omitempty" example:""`
}
