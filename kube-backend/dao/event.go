package dao

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/wonderivan/logger"
	"kube-backend/db"
	"kube-backend/model"
	"time"
)

var Event event

type event struct {
}

// 定义event列表结构体
type Events struct {
	Items []*model.Event `json:"items"`
	Total int            `json:"total"`
}

func (*event) GetList(name, cluster string, page, limit int) (*Events, error) {
	//定义分页数据起始位置
	startSet := (page - 1) * limit
	//定义数据库查询的返回内容
	var (
		eventList = make([]*model.Event, 0)
		total     = 0
	)
	tx := db.GORM.
		Model(&model.Event{}).
		Where("name like ? and cluster = ?", "%"+name+"%", cluster).
		Count(&total).Limit(limit).Offset(startSet).Order("id desc").
		Find(&eventList)
	if tx.Error != nil {
		logger.Error(fmt.Sprintf("查询event列表失败，%v\n", tx.Error))
		return nil, errors.New(fmt.Sprintf("查询event列表失败，%v\n", tx.Error))
	}
	return &Events{
		Items: eventList,
		Total: total,
	}, nil
}

// 新增
func (*event) Add(event *model.Event) error {
	tx := db.GORM.Create(&event)
	if tx.Error != nil {
		logger.Error(fmt.Sprintf("添加event列表失败，%v\n", tx.Error))
		return errors.New(fmt.Sprintf("添加event列表失败，%v\n", tx.Error))
	}
	return nil
}

// 查询单个event.因为新增的可能有重复的，所以每次add后去查询下单个的event是否存在，存在就不加了.
func (*event) HasEvent(name, kind, namespace, reason string, eventTime time.Time, cluster string) (*model.Event, bool, error) {
	data := &model.Event{}
	tx := db.GORM.Where("name = ? and kind = ? and namespace = ? and reason = ? and event_time = ? and cluster = ?",
		name, kind, namespace, reason, eventTime, cluster).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if tx.Error != nil {
		logger.Error(fmt.Sprintf("查询单个event列表失败，%v\n", tx.Error))
		return nil, false, errors.New(fmt.Sprintf("查询单个event列表失败，%v\n", tx.Error))
	}
	return data, true, nil
}
