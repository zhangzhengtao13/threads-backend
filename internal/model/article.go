package model

import (
	"threads-service/pkg/app"

	"github.com/jinzhu/gorm"
)

type Article struct {
	*Model
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	State         uint8  `json:"state"`
}

func (a Article) TableName() string {
	return "blog_article"
}

type ArticleSwagger struct {
	List  []*Article
	Pager *app.Paper
}

func (a Article) CreateArticle(db *gorm.DB) (*Article, error) {
	err := db.Create(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (a Article) DeleteArticle(db *gorm.DB) error {
	err := db.Where("id=? and deleted=?", a.ID, 0).Delete(&a).Error
	if err != nil {
		return err
	}
	return nil
}

func (a Article) UpdateArticle(db *gorm.DB, new interface{}) error {
	err := db.Table(a.TableName()).Where("id=? and deleted=?", a.ID, 0).Updates(new).Error
	if err != nil {
		return err
	}
	return nil
}

func (a Article) GetAnArticle(db *gorm.DB) (Article, error) {
	var article Article
	db = db.Table(a.TableName()).Where("id=? and deleted=? and state=?", a.ID, 0, a.State)
	err := db.Take(&article).Error // 确保err是由Take产生的, 不是其他产生的
	if err != nil && err != gorm.ErrRecordNotFound {
		return article, err
	}
	return article, nil
}

// func (a Article) List(db *gorm.DB, pageOffset, pageSize int) ([]*Article, error) {
// 	var articles []*Article
// 	if pageOffset >= 0 && pageSize > 0 {
// 		db.Offset(pageOffset).Limit(pageSize)
// 	}
// 	err := db.Table(a.TableName()).Where("deleted=?", 0).Find(&articles).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return articles, nil
// }

type ArticleRow struct {
	ID            uint32
	TagID         uint32
	TagName       string
	Title         string
	Description   string
	CoverImageURL string
	Content       string
}

func (a Article) ListByTag(db *gorm.DB, tagID uint32, pageOffset, pageSize int) ([]*ArticleRow, error) {
	fields := []string{"ar.id as article_id", "ar.title as article_title", "ar.desc as article_desc", "ar.cover_image_url", "ar.content"}
	fields = append(fields, []string{"t.id as tag_id", "t.name as tag_name"}...)

	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	rows, err := db.Select(fields).Table(ArticleTag{}.TableName()+" as at").Joins("left outer join "+Tag{}.TableName()+" as t on at.tag_id=t.id").Joins("left join "+Article{}.TableName()+" as ar on at.article_id=ar.id").Where("at.`tag_id`=? and ar.state=? and ar.deleted=?", tagID, a.State, 0).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var articles []*ArticleRow
	for rows.Next() {
		ar := &ArticleRow{}
		err = rows.Scan(&ar.ID, &ar.Title, &ar.Description, &ar.CoverImageURL, &ar.Content, &ar.TagID, &ar.TagName)
		if err != nil {
			return nil, err
		}
		articles = append(articles, ar)
	}
	return articles, nil
}

func (a Article) CountByTagID(db *gorm.DB, tagID uint32) (int, error) {
	var count int
	err := db.Table(ArticleTag{}.TableName()+" as at").Joins("left outer join "+Tag{}.TableName()+" as t on at.tag_id=t.id").Joins("left outer join "+Article{}.TableName()+" as ar on at.article_id=ar.id").Where("at.tag_id=? and ar.state=? and ar.deleted=?", tagID, a.State, 0).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
