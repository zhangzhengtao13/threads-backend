package service

import (
	"threads-service/internal/dao"
	"threads-service/internal/model"
	"threads-service/pkg/app"
)

/*-----------------------对业务参数的数据抽象--------------------*/

type ArticleRequest struct {
	ID    uint32 `form:"id" binding:"required,gte=1"`
	State uint8  `form:"state, default=1" binding:"oneof=0 1"`
}

type ArticleListRequest struct {
	TagID uint32 `form:"tag_id" binding:"gte=1"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type CreateArticleRequest struct {
	TagID         uint32 `form:"tag_id" binding:"required,gte=1"`
	Title         string `form:"title" binding:"required,min=2,max=100"`
	Desc          string `form:"desc" binding:"required,min=2,max=255"`
	Content       string `form:"content" binding:"required,min=2,max=4294967295"`
	CoverImageUrl string `form:"cover_image_url" binding:"required,url"`
	CreatedBy     string `form:"created_by" binding:"required,min=2,max=100"`
	State         uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type UpdateArticleRequest struct {
	ID            uint32 `form:"id" binding:"required,gte=1"`
	TagID         uint32 `form:"tag_id" binding:"required,gte=1"`
	Title         string `form:"title" binding:"min=2,max=100"`
	Desc          string `form:"desc" binding:"min=2,max=255"`
	Content       string `form:"content" binding:"min=2,max=4294967295"`
	CoverImageUrl string `form:"cover_image_url" binding:"url"`
	ModifiedBy    string `form:"modified_by" binding:"required,min=2,max=100"`
	State         uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type DeleteArticleRequest struct {
	ID uint32 `form:"id" binding:"required,gte=1"`
}

// 对主角数据的抽象
type Article struct {
	ID            uint32     `json:"id"`
	Title         string     `json:"title"`
	Desc          string     `json:"desc"`
	Content       string     `json:"content"`
	CoverImageUrl string     `json:"cover_image_url"`
	State         uint8      `json:"state"`
	Tag           *model.Tag `json:"tag"`
}

// 过程抽象

func (svc *Service) GetArticle(param *ArticleRequest) (*Article, error) {
	article, err := svc.dao.GetAnArticle(param.ID, param.State)
	if err != nil {
		return nil, err
	}
	at, err := svc.dao.GetArticleTagByAID(article.ID)
	if err != nil {
		return nil, err
	}
	tag, err := svc.dao.GetTag(at.TagID, model.STATE_ON)
	if err != nil {
		return nil, err
	}
	return &Article{
		ID:            article.ID,
		Title:         article.Title,
		Desc:          article.Desc,
		Content:       article.Content,
		CoverImageUrl: article.CoverImageUrl,
		State:         article.State,
		Tag:           &tag,
	}, nil
}

func (svc *Service) GetArticleList(param *ArticleListRequest, pager *app.Paper) ([]*Article, int, error) {
	count, err := svc.dao.CountArticleListByTagID(param.TagID, param.State)
	if err != nil {
		return nil, 0, err
	}
	ars, err := svc.dao.GetArticleListByTagID(param.TagID, param.State, pager.Page, pager.PageSize)
	if err != nil {
		return nil, 0, err
	}
	var list []*Article
	for _, ar := range ars {
		list = append(list, &Article{
			ID:            ar.ID,
			Title:         ar.Title,
			Desc:          ar.Description,
			Content:       ar.Content,
			CoverImageUrl: ar.CoverImageURL,
			Tag:           &model.Tag{Model: &model.Model{ID: ar.TagID}, Name: ar.TagName},
		})
	}
	return list, count, nil
}

func (svc *Service) CreateArticle(param *CreateArticleRequest) error {
	ar, err := svc.dao.CreatedArticle(&dao.Article{
		Title:         param.Title,
		Description:   param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		CreatedBy:     param.CreatedBy,
		State:         param.State,
	})
	if err != nil {
		return err
	}
	err = svc.dao.CreateArticleTag(ar.ID, param.TagID, param.CreatedBy)
	if err != nil {
		return err
	}
	return nil
}

func (svc *Service) UpdateArticle(param *UpdateArticleRequest) error {
	err := svc.dao.UpdateArticle(&dao.Article{
		ID:            param.ID,
		Title:         param.Title,
		Description:   param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		State:         param.State,
		ModifiedBy:    param.ModifiedBy,
	})
	if err != nil {
		return err
	}
	err = svc.dao.UpdateArticleTag(param.ID, param.TagID, param.ModifiedBy)
	if err != nil {
		return err
	}
	return nil
}

func (svc *Service) DeleteArticle(param *DeleteArticleRequest) error {
	err := svc.dao.DeleteArticle(param.ID)
	if err != nil {
		return err
	}
	err = svc.dao.DeleteArticleTag(param.ID)
	if err != nil {
		return err
	}
	return nil
}
