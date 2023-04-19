package v1

import (
	"threads-service/global"
	"threads-service/internal/service"
	"threads-service/pkg/app"
	"threads-service/pkg/convert"
	"threads-service/pkg/errcode"

	"github.com/gin-gonic/gin"
)

type Tag struct {
}

func NewTag() Tag {
	return Tag{}
}

// @Summary	获取多个标签
// @Produce	json
// @Param		name		query		string		fasle	"标签名称"	maxlength(100)
// @Param		state		query		int			fasle	"状态"	Enums(0, 1)	default(1)
// @Param		page		query		int			false	"页码"
// @Param		page_size	query		int			false	"每页数量"
// @Success	200			{object}	model.TagSwagger	"成功"
// @Failure	400			{boject}	errcode.Err	"请求错误"
// @Failure	500			{object}	errcode.Err	"内部错误"
// @Router		/api/v1/tags [get]
func (t Tag) List(c *gin.Context) {
	param := service.TagListRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	pager := app.Paper{
		Page:     app.GetPage(c),
		PageSize: app.GetPageSize(c),
	}
	totalRows, err := svc.CountTag(&service.CountTagRequest{Name: param.Name, State: param.State})
	if err != nil {
		global.Logger.Errf(c, "svc.CountTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorCountTagFail)
		return
	}
	tags, err := svc.GetTagList(&param, &pager)
	if err != nil {
		global.Logger.Errf(c, "svc.GetTagList err:%v", err)
		response.ToErrorResponse(errcode.ErrorGetTagListFail)
		return
	}
	response.ToResponseList(tags, totalRows)
}

// @Summary	新增标签
// @Produce	json
// @Param		name		body		string		true	"标签名称"	minlength(3)	maxlength(100)
// @Param		state		body		int			false	"状态"	Enums(0, 1)		default(1)
// @Param		created_by	body		string		false	"创建者"	minlength(3)	maxlength(100)

// @Success	200			{object}	model.TagSwagger	"成功"
// @Failure	400			{object}	errcode.Err	"请求错误"
// @Failure	500			{object}	errcode.Err	"内部错误"
// @Router		/api/v1/tags [post]
func (t Tag) Create(c *gin.Context) {
	param := service.CreateTagRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	err := svc.CreateTag(&param)
	if err != nil {
		global.Logger.Errf(c, "svc.CreateTag err:%v", err)
		response.ToErrorResponse(errcode.ErrorCraeteTagFail)
		return
	}
	response.ToResponse(gin.H{})
}

//	@Summary	更新标签
//	@Produce	json
//	@Param		id			path	int		true	"标签ID"
//	@Param		name		body	string	false	"标签名称"	minlength(3)	maxlength(100)
//	@Param		state		body	int		false	"状态"	Enums(0, 1)		default(1)
//	@Param		modified_by	body	string	true	"修改者"	minlength(3)	maxlength(100)
//
// Success 200 {array} model.TagSwagger "成功"
// Failure 400 {object} errcode.Err "请求错误"
// Failure 500 {object} errcode.Err "内部错误"
// Router /api/v1/tags/{id} [put]
func (t Tag) Update(c *gin.Context) {
	param := service.UpdateTagRequest{
		ID: convert.StrTo(c.Param("id")).MustUInt32(),
	}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errf(c.Request.Context(), "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	err := svc.UpdateTag(&param)
	if err != nil {
		global.Logger.Errf(c, "svc.UpdateTag err:%v", err)
		response.ToErrorResponse(errcode.ErrorUpdateTagFail)
		return
	}
	response.ToResponse(gin.H{})
}

//	@Summary	删除标签
//	@Produce	json
//	@Param		id	path	int	true	"标签ID"
//
// Success 200 {string} string "成功"
// Failure 400 {object} errcode.Err "请求错误"
// Failure 500 {object} errcode.Err "服务器错误"
//
//	@Router		/api/v1/tags/{id} [delete]
func (t Tag) DELETE(c *gin.Context) {
	param := service.DeleteTagRequest{
		ID: convert.StrTo(c.Param("id")).MustUInt32(),
	}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errf(c, "app.BindAndValid err: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	err := svc.DeleteTag(&param)
	if err != nil {
		global.Logger.Errf(c, "svc.DeletTag err:%v", err)
		response.ToErrorResponse(errcode.ErrorDeleteTagFail)
		return
	}
	response.ToResponse(gin.H{})
}
