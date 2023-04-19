package global

import (
	"github.com/jinzhu/gorm"
)

var (
	// 数据库
	DBEngine *gorm.DB // main.go文件中获取值
)
