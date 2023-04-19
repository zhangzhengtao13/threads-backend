package errcode

var (
	// 文件上传业务
	ERROR_UPLOAD_FILE_FAIL = NewError(20030001, "上传文件失败")
)
