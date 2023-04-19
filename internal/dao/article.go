package dao

import (
	"threads-service/internal/model"
	"threads-service/pkg/app"
)

type Article struct {
	ID            uint32 `json:"id"`
	TagID         uint32 `json:"tag_id"`
	Title         string `json:"title"`
	Description   string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
	State         uint8  `json:"state"`
}

func (d *Dao) CreatedArticle(param *Article) (*model.Article, error) {
	article := model.Article{
		Title:         param.Title,
		Desc:          param.Description,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		State:         param.State,
		Model:         &model.Model{CreatedBy: param.CreatedBy},
	}
	return article.CreateArticle(d.engine)
}

func (d *Dao) UpdateArticle(param *Article) error {
	article := model.Article{
		Model: &model.Model{ID: param.ID},
	}
	value := map[string]interface{}{
		"modified_by": param.ModifiedBy,
		"state":       param.State,
	}
	if param.Title != "" {
		value["title"] = param.Title
	}
	if param.CoverImageUrl != "" {
		value["cover_image_url"] = param.CoverImageUrl
	}
	if param.Description != "" {
		value["desc"] = param.Description
	}
	if param.Content != "" {
		value["content"] = param.Content
	}
	return article.UpdateArticle(d.engine, value)
}

func (d *Dao) GetAnArticle(id uint32, state uint8) (model.Article, error) {
	article := model.Article{Model: &model.Model{ID: id}, State: state}
	return article.GetAnArticle(d.engine)
}

func (d *Dao) DeleteArticle(id uint32) error {
	ar := model.Article{Model: &model.Model{ID: id}}
	return ar.DeleteArticle(d.engine)
}

func (d *Dao) CountArticleListByTagID(id uint32, state uint8) (int, error) {
	ar := model.Article{State: state}
	return ar.CountByTagID(d.engine, id)
}

func (d *Dao) GetArticleListByTagID(id uint32, state uint8, page, pageSize int) ([]*model.ArticleRow, error) {
	ar := model.Article{State: state}
	return ar.ListByTag(d.engine, id, app.GetPageOffset(page, pageSize), pageSize)
}
