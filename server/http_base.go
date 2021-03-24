package server

//响应的消息
type ResponseData struct {
	Stat    int         `json:"stat"`
	Data    interface{} `json:"data"`
}

func BuildResponse(stat int, data interface{}) *ResponseData {
	return &ResponseData{
		Stat:    stat,
		Data:    data,
	}
}