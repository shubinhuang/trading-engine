package errcode

// 各种错误码
var (
	OK             = NewErrCode(20, "20: OK")
	UNDEFINE       = NewErrCode(23, "23: Undefine")
	ChanFull       = NewErrCode(40, "40: Order Channel is full")
	EngineExist    = NewErrCode(41, "41: Trade engine already exist")
	EngineNotFound = NewErrCode(42, "42: Trade engine not found")
	OrderExist     = NewErrCode(43, "43: Order already exist")
	OrderNotFound  = NewErrCode(44, "44: Order not found")
	FORMATERROR    = NewErrCode(45, "45: Request format error")
)

// 撮合服务名称
var ServiceName = "servicename"
