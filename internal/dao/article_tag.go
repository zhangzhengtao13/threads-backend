package dao

import "threads-service/internal/model"

func (d *Dao) GetArticleTagByAID(aID uint32) (model.ArticleTag, error) {
	at := model.ArticleTag{ArticleID: aID}
	return at.GetByAID(d.engine)
}

func (d *Dao) GetArticleTagListByTID(tagID uint32) ([]*model.ArticleTag, error) {
	at := model.ArticleTag{TagID: tagID}
	return at.ListByTID(d.engine)
}

func (d *Dao) GetArticleTagByAIDs(aids []uint32) ([]*model.ArticleTag, error) {
	at := model.ArticleTag{}
	return at.ListByAIDs(d.engine, aids)
}

func (d *Dao) CreateArticleTag(aID, tID uint32, creator string) error {
	at := model.ArticleTag{
		Model:     &model.Model{CreatedBy: creator},
		ArticleID: aID,
		TagID:     tID,
	}
	return at.Create(d.engine)
}

func (d *Dao) UpdateArticleTag(aID, tID uint32, modifiedBy string) error {
	at := model.ArticleTag{
		ArticleID: aID,
	}
	vs := map[string]interface{}{
		"article_id":  aID,
		"tag_id":      tID,
		"modified_by": modifiedBy,
	}
	return at.UpdateOne(d.engine, vs)
}

func (d *Dao) DeleteArticleTag(aid uint32) error {
	at := model.ArticleTag{
		ArticleID: aid,
	}
	return at.DeleteOne(d.engine)
}
