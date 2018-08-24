package controller

// Result represents HTTP response body.
type Result struct {
	Code int
	Msg  string
	Data interface{}
}

// NewResult creates a result with Code=0, Msg="", Data=nil.
func NewResult() *Result {
	return &Result{
		Code: Success,
		Msg:  SuccessMsg,
		Data: nil,
	}
}

const (
	Success            = 0
	SuccessMsg         = "成功"
	ParamsError        = -1001
	ParamsErrorMsg     = "参数错误"
	CreateError        = -1002
	CreateErrorMsg     = "创建失败"
	UpdateError        = -1003
	UpdateErrorMsg     = "修改失败"
	DeleteError        = -1004
	DeleteErrorMsg     = "删除失败"
	OperateRepeat      = -1005
	OperateRepeatMsg   = "操作重复"
	OperateError       = -1006
	OperateErrorMsg    = "操作失败"
	ParamsDataError    = -1007
	ParamsDataErrorMsg = "参数数据错误"
	ServiceTimeOut     = -1008
	ServiceTimeOutMsg  = "服务当前繁忙"
	ServiceError       = -1009
	ServiceErrorMsg    = "服务异常"
)

const (
	ForbiddenError = -1403
	NoPowerUpdate  = "没权限修改"
)
