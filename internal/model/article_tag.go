package model

import "github.com/jinzhu/gorm"

type ArticleTag struct {
	*Model
	TagID     uint32 `json:"tag_id"`
	ArticleID uint32 `json:"article_id"`
}

func (at ArticleTag) TableName() string {
	return "blog_article_tag"
}

func (at ArticleTag) GetByAID(db *gorm.DB) (ArticleTag, error) {
	var aT ArticleTag
	err := db.Table(at.TableName()).Where("article_id=? and deleted=?", at.ArticleID, 0).Take(&aT).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return aT, err
	}
	return aT, nil
}

func (at ArticleTag) ListByTID(db *gorm.DB) ([]*ArticleTag, error) {
	var ats []*ArticleTag
	err := db.Table(at.TableName()).Where("tag_id=? and deleted=?", at.TagID, 0).Find(&ats).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return ats, nil
}

func (at ArticleTag) ListByAIDs(db *gorm.DB, aids []uint32) ([]*ArticleTag, error) {
	var ats []*ArticleTag
	err := db.Table(at.TableName()).Where("article_id in ? and deleted=?", aids, 0).Find(&ats).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return ats, nil
}

func (at ArticleTag) Create(db *gorm.DB) error {
	err := db.Create(&at).Error
	if err != nil {
		return err
	}
	return nil
}

func (at ArticleTag) UpdateOne(db *gorm.DB, vs interface{}) error {
	err := db.Table(at.TableName()).Where("article_id=? and deleted=?", at.ArticleID, 0).Updates(vs).Error
	if err != nil {
		return err
	}
	return nil
}

func (at ArticleTag) DeleteMany(db *gorm.DB) (err error) {
	err = db.Table(at.TableName()).Where("id=? and deleted=?", at.ID, 0).Delete(&at).Error
	if err != nil {
		return
	}
	return nil
}

func (at ArticleTag) DeleteOne(db *gorm.DB) error {
	err := db.Table(at.TableName()).Where("article_id=> and deleted=?", at.ArticleID, 0).Limit(1).Delete(&at).Error
	if err != nil {
		return err
	}
	return nil
}
