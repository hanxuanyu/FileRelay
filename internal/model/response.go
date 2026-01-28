package model

// Response 统一响应模型
type Response struct {
	Code int         `json:"code" example:"200"`
	Msg  string      `json:"msg" example:"success"`
	Data interface{} `json:"data"`
}

// 错误码定义
const (
	CodeSuccess         = 200
	CodeBadRequest      = 400
	CodeUnauthorized    = 401
	CodeForbidden       = 403
	CodeNotFound        = 404
	CodeGone            = 410
	CodeInternalError   = 500
	CodeTooManyRequests = 429
)

// NewResponse 创建响应
func NewResponse(code int, msg string, data interface{}) Response {
	return Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// SuccessResponse 成功响应
func SuccessResponse(data interface{}) Response {
	return Response{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	}
}

// ErrorResponse 错误响应
func ErrorResponse(code int, msg string) Response {
	return Response{
		Code: code,
		Msg:  msg,
	}
}
