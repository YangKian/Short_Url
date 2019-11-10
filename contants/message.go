package contants

var CodeMessage = map[int]string{
	SUCCESS: "OK",
	ERROR:   "Failed",

	NOTFOUND:     "找不到资源",
	UNAUTHORIZED: "未授权",

	DBERROER:               "数据库错误",
	CREATE_SHORT_URL_ERROR: "创建短链接失败",
	REQUEST_ERROR:          "请求错误",
	EMPTYREQUESTBODY:       "请求体不能为空",
}

func MsgGeter(code int) string {
	if msg, ok := CodeMessage[code]; ok {
		return msg
	}

	return CodeMessage[ERROR]
}
