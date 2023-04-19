package global

import (
	"threads-service/pkg/logger"
	"threads-service/pkg/setting"
)

var (
	// 配置
	ServerSetting   *setting.ServerSettings
	AppSetting      *setting.AppSettings
	DataBaseSetting *setting.DataBaseSettings

	// 日志
	Logger *logger.Logger

	// auth
	JwtSetting *setting.JWTSettings

	// 邮箱
	EmailSetting *setting.EmailSettings
)
