package errorcode

type Value int

const (
	Null   Value = 0   // Null is 无效错误
	API    Value = 501 // API is api错误码
	Verify Value = 502 // Verify is 验证错误码
	Panic  Value = 599 // Panic is 异常错误码
)
