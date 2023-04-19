package routers

import (
	"threads-service/global"
	"threads-service/internal/service"
	"threads-service/pkg/app"
	"threads-service/pkg/convert"
	"threads-service/pkg/errcode"
	"threads-service/pkg/upload"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	response := app.NewResponse(c)
	file, fileHeader, err := c.Request.FormFile("file")
	fileType := convert.StrTo(c.PostForm("type")).MustInt()
	if err != nil {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}
	if fileHeader == nil || fileType <= 0 {
		response.ToErrorResponse(errcode.InvalidParams)
	}
	svc := service.New(c.Request.Context())
	fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, fileHeader)
	if err != nil {
		global.Logger.Errf(c.Request.Context(), "svc.UploadFile err: %v", err)
		response.ToErrorResponse(errcode.ERROR_UPLOAD_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	response.ToResponse(gin.H{
		"file_access_url": fileInfo.AccessUrl,
	})
}
