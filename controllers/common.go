package controllers

var errorCode map[int]string

type resultStruct struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func GetMessage(code int, data interface{}) resultStruct {
	//if _, exists := errorCode[code]; !exists {
	//	code = 998
	//}

	return resultStruct{Code: code, Msg: errorCode[code], Data: data}
}
