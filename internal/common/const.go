package common

const (
	SuccessCode = 0
	ErrorCode   = 1
)

var msgList = map[int32]string{
	SuccessCode: "success",
	ErrorCode:   "error",
}

func ResponseMsg(code int32) string {
	if msg, ok := msgList[code]; ok {
		return msg
	}
	return "unknown msg"
}
