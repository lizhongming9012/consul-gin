package e

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",

	INVALID_PARSE_FORM:              "解析绑定表单错误",
	INVALID_PARAMS_VERIFY:           "参数校验错误",
	ERROR_EXIST_LICENCE:             "已存在该LICENCE",
	ERROR_EXIST_LICENCE_FAIL:        "获取已存在LICENCE失败",
	ERROR_NOT_EXIST_LICENCE:         "该LICENCE不存在",
	ERROE_VERSION_LOW:               "版本过低，请更新新版本",
	ERROR_USERNAME_PASSWORD:         "用户名密码不正确",
	ERROR_AUTH_CHECK_TOKEN_FAIL:     "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:  "Token已超时,请重新登录",
	ERROR_AUTH_TOKEN:                "Token生成失败",
	ERROR_AUTH:                      "Token错误",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "保存图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "检查图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "校验图片错误，图片格式或大小有问题",

	ERROR_UPLOAD_SAVE_FILE_FAIL:    "保存文件失败",
	ERROR_UPLOAD_CHECK_FILE_FAIL:   "检查文件失败",
	ERROR_UPLOAD_CHECK_FILE_FORMAT: "校验文件错误，文件格式不正确",
	ERROR_UPLOAD_CHECK_FILE_SIZE:   "校验文件错误，文件大小超出限制",

	ERROR_USERNAME_EXIST: "用户名已存在",

	ERROR_GET_DEPARTMENT_FAIL: "获取部门列表失败",
	ERROR_GET_USER_FAIL:       "获取部门用户列表失败",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
