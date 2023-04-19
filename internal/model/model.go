package model

import (
	"fmt"
	"threads-service/global"
	"threads-service/pkg/setting"
	"time"

	otgorm "github.com/eddycjy/opentracing-gorm"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	STATE_ON  = 1
	STATE_OFF = 0
)

type Model struct {
	ID         uint32 `gorm:"primary key" josn:"id,omitempty"`
	CreatedBy  string `json:"created_by,omitempty"`
	ModifiedBy string `json:"modified_by,omitempty"`
	CreatedOn  uint32 `json:"created_on,omitempty"`
	ModifiedOn uint32 `json:"modified_on,omitempty"`
	DeletedOn  uint32 `json:"deleted_on,omitempty"`
	Deleted    uint8  `json:"deleted,omitempty"`
}

func NewDBEngine(databaseSetting *setting.DataBaseSettings) (*gorm.DB, error) {
	db, err := gorm.Open(databaseSetting.DBType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local", databaseSetting.UserName, databaseSetting.Password, databaseSetting.Host, databaseSetting.DBName, databaseSetting.Charset, databaseSetting.ParseTime))
	if err != nil {
		return nil, err
	}
	if global.ServerSetting.RunMode == "debug" {
		db.LogMode(true)
	}
	db.SingularTable(true)

	db.Callback().Create().Replace("gorm:update_time_stamp", updateStampForUpdateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallBack)

	db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	db.DB().SetMaxOpenConns(databaseSetting.MamOpenConns)

	otgorm.AddGormCallbacks(db)
	return db, nil
}

/*回调函数: 处理公共字段*/

// 新增回调
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				_ = createTimeField.Set(nowTime)
			}
		}
		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				_ = modifyTimeField.Set(nowTime)
			}
		}
	}
}

func updateStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		_ = scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

// 判断是否存在DeletedOn和IsDel字段。若存在，则执行UPDATE操作进行软删除（修改DeletedOn和IsDel的值），否则执行DELETE操作进行硬删除.
// 调用scope.QuotedTableName方法获取当前引用的表名，并调用一系列方法对SQL语句的组成部分进行处理和转移。
// 在完成一些所需参数设置后，调用scope.CombinedConditionSql方法完成SQL语句的组装
func deleteCallBack(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}
		deletedOnField, hasDeleteOnField := scope.FieldByName("DeletedOn")
		isDeletedField, hasIsDelField := scope.FieldByName("Deleted")
		if !scope.Search.Unscoped && hasDeleteOnField && hasIsDelField {
			now := time.Now().Unix()
			scope.Raw(fmt.Sprintf("UPDATE %v set %v=%v, %v=%v%v%v", scope.QuotedTableName(), scope.Quote(deletedOnField.DBName), scope.AddToVars(now),
				scope.Quote(isDeletedField.DBName),
				scope.AddToVars(1),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption))).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETED from %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}

	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
