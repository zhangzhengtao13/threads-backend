package v1

import (
	"net/http"
	"threads-service/global"
	"threads-service/internal/model"
	"threads-service/internal/service"
	"threads-service/pkg/app"
	"threads-service/pkg/convert"
	"threads-service/pkg/errcode"

	"github.com/gin-gonic/gin"
)

type Article struct {
}

func NewArticle() Article {
	return Article{}
}

//	@Summary	获取单个文章
//	@Produce	json
//	@Param		id	query		int			false	"文章ID"
//	@Success	200	{object}	model.ArticleSwagger	"成功"
//
// Failure 400 {object} errcode.Err 请求错误
// Failure 500 {object} errcode.Err 服务器错误
//
//	@Router		/api/v1/articles/{id} [get]
func (a Article) Get(c *gin.Context) {
	param := service.ArticleRequest{
		ID:    convert.StrTo(c.Param("id")).MustUInt32(),
		State: model.STATE_ON,
	}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errf(c.Request.Context(), "app.BindAndValid err:%v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	ar, err := svc.GetArticle(&param)
	if err != nil {
		global.Logger.Errf(c.Request.Context(), "svc.GetArticle err:%v", err)
		response.ToErrorResponse(errcode.ErrorGetArticleFail)
		return
	}
	response.ToResponse(ar)
}

// @Summary	获取多个文章
// @Produce	json
// @Param		title		query		string		false	"文章标题"	maxlength(30)
// @Param		state		query		int			false	"文章状态"	Enums(0, 1)	default(1)
// @Param		page		query		int			false	"页码"
// @Param		page_size	query		int			false	"每页数量"
// @Success	200			{object}	model.ArticleSwagger	"成功"
// @Failure	400			{object}	errcode.Err	"请求错误"
// @Failure	500			{object}	errcode.Err	"服务器错误"
// @Router		/api/v1/articles [get]
func (a Article) List(c *gin.Context) {
	param := service.ArticleListRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errf(c.Request.Context(), "app.BindAndValid err:%v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	paper := app.Paper{Page: app.GetPage(c), PageSize: app.GetPageSize(c)}
	ars, rows, err := svc.GetArticleList(&param, &paper)
	if err != nil {
		global.Logger.Errf(c.Request.Context(), "svc.GetArticleList err:%v", err)
		response.ToErrorResponse(errcode.ErrorGetArticlesFail)
		return
	}
	response.ToResponseList(ars, rows)
}

// @Summary	创建文章
// @Produce	json
// @Param		title				body		string		true	"文章标题"	minlength(1)	maxlength(30)
// @Param		state				body		int			false	"状态"	Enums(0, 1)		default(1)
// @Param		desc				body		string		false	"描述"	maxlength(50)
// @Param		content				body		string		true	"内容"	minlength(10)	maxlength(5000)
// @Param		cover_image_urlbody	body		string		false	"封面图片地址"
// @Param		created_by			body		string		false	"发表人"	maxlength(10)
// @Success	200					{object}	model.ArticleSwagger	"成功"
// @Failure	400					{object}	errcode.Err	"请求错误"
// @Failure	500					{object}	errcode.Err	"服务器错误"
// @Router		/api/v1/articles [post]
func (a Article) Create(c *gin.Context) {
	param := service.CreateArticleRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errf(c.Request.Context(), "app.BindAndValid err: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	err := svc.CreateArticle(&param)
	if err != nil {
		global.Logger.Errf(c.Request.Context(), "svc.CraeteArticle err:%v", err)
		response.ToErrorResponse(errcode.ErrorCreateArticleFail)
		return
	}
	response.ToResponse(gin.H{"status": http.StatusOK, "message": "创建文章成功"})
}

// @Summary	修改文章
// @Produce	json
// @Param		id					path		int			false	"文章ID"
// @Param		title				body		string		false	"文章标题"	minlength(1)	maxlength(30)
// @Param		state				body		int			false	"状态"	Enums(0,1)		default(1)
// @Param		desc				body		string		false	"描述"	maxlength(50)
// @Param		content				body		string		false	"正文"	minlength(10)	maxlength(5000)
// @Param		cover_image_urlbody	body		string		false	"封面图片地址"
// @Param		modified_by			body		string		true	"发表人"	maxlength(10)
// @Success	200					{object}	model.ArticleSwagger	"成功"
// @Failure	400					{object}	errcode.Err	"请求错误"
// @Failure	500					{object}	errcode.Err	"服务器错误"
// @Router		/api/v1/articles{id} [put]
func (a Article) Update(c *gin.Context) {
	param := service.UpdateArticleRequest{ID: convert.StrTo(c.Param("id")).MustUInt32()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errf(c.Request.Context(), "app.BindAndValid err: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	err := svc.UpdateArticle(&param)
	if err != nil {
		global.Logger.Errf(c.Request.Context(), "svc.UpdateArticle err:%v", err)
		response.ToErrorResponse(errcode.ErrorUpdateArticleFail)
		return
	}
	response.ToResponse(gin.H{})
}

//	@Summary	删除文章
//	@Produce	json
//	@Param		id	path	int	true	"标签ID"
//
// Success 200 {string} string "成功"
// Failure 400 {object} errcode.Err "请求错误"
// Failure 500 {object} errcode.Err "服务器错误"
//
//	@Router		/api/v1/articles/{id} [delete]
func (a Article) Delete(c *gin.Context) {
	para := service.DeleteArticleRequest{ID: convert.StrTo(c.Param("id")).MustUInt32()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &para)
	if !valid {
		global.Logger.Errf(c, "app.BindAndValid err: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	err := svc.DeleteArticle(&para)
	if err != nil {
		global.Logger.Errf(c, "svc.DeleteArticle err: %v", err)
		response.ToErrorResponse(errcode.ErrorDeleteArticleFail)
		return
	}
	response.ToResponse(gin.H{})
}
