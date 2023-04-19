package model

import "github.com/jinzhu/gorm"

type Auth struct {
	ID        uint32
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

func (a Auth) TableName() string {
	return "blog_auth"
}

func (a Auth) Get(db *gorm.DB) (Auth, error) {
	var auth Auth
	db = db.Where("app_key=? and app_secret=? and deleted=?", a.AppKey, a.AppSecret, 0)
	err := db.Take(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return auth, err
	}
	return auth, nil
}
