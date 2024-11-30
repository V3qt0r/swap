package api

type Response struct {
	Code 				int		 		`json:"code" yaml:"code" example:"500"`
	Message 			string    		`json:"message" yaml:"message"`
	Details             interface{}		`json:"details" yaml:"details"`
}


func NewResponse(Code int, Message string, Details interface{}) *Response {
	return &Response{
		Code: 		Code,
		Message: 	Message,
		Details: 	Details,
	}
}