package model

import (
	"threads-service/pkg/app"

	"github.com/jinzhu/gorm"
)

type Tag struct {
	*Model
	Name  string `json:"name"`
	State uint8  `json:"state"`
}

// 返回数据库表名
func (t Tag) TableName() string {
	return "blog_tag"

}

type TagSwagger struct {
	List  []*Tag
	Pager *app.Paper
}

func (t Tag) Count(db *gorm.DB) (int, error) {
	var count int
	if t.Name != "" {
		db = db.Where("name=?", t.Name)
	}
	db = db.Where("state=?", t.State)
	if err := db.Model(&t).Where("deleted=?", 0).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (t Tag) List(db *gorm.DB, pageOffset, pageSize int) ([]*Tag, error) {
	var tags []*Tag
	var err error
	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	err = db.Where("name=?", t.TableName()).Where("deleted=?", 0).Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (t Tag) Create(db *gorm.DB) error {
	return db.Create(&t).Error
}

func (t Tag) Update(db *gorm.DB, values interface{}) error {
	err := db.Model(&t).Updates(values).Where("id=? AND deleted=?", t.ID, 0).Error
	if err != nil {
		return err
	}
	return nil
}

func (t Tag) Delete(db *gorm.DB) error {
	return db.Where("id=? AND deleted=?", t.ID, 0).Delete(&t).Error
}

func (t Tag) Get(db *gorm.DB) (Tag, error) {
	var tag Tag
	err := db.Table(t.TableName()).Where("id=? and deleted=?", t.ID, 0).Take(&tag).Error
	if err != nil {
		return tag, err
	}
	return tag, nil
}
