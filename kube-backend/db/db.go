package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wonderivan/logger"
	"kube-backend/config"
	"kube-backend/model"
)

var (
	isInit bool //是否已经初始化
	GORM   *gorm.DB
	err    error
)

// db初始化
func Init() {
	if isInit {
		return
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DbUser,
		config.DbPwd,
		config.DbHost,
		config.DbPort,
		config.DbName)
	GORM, err = gorm.Open(config.DbType, dsn)
	GORM.LogMode(config.LogMode)
	GORM.DB().SetMaxIdleConns(config.MaxIdleConns)
	GORM.DB().SetMaxOpenConns(config.MaxOpenConns)
	GORM.DB().SetConnMaxLifetime(config.MaxLifeTime)
	isInit = true
	GORM.AutoMigrate(model.Event{}, model.Chart{})
	logger.Info("连接数据库成功")
}

func Close() error {
	logger.Info("关闭数据库连接")
	return GORM.Close()
}
