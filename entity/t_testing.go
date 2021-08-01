package entity

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Testing struct {
	ID        uint   `gorm:"primaryKey"`
	EnvId     uint   `gorm:"comment:环境"`
	ApiId     uint   `gorm:"comment:所属 api"`
	UserId    uint   `gorm:"comment:所属用户"`
	Name      string `gorm:"not null;comment:名称"`
	Project   string `gorm:"not null;comment:项目"`
	Data      string `gorm:"type:text;not null;comment:测试数据"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Testing) TableName() string {
	return "apitesting_testing"
}

func (Testing) DataInit() error {
	return nil
}

func (Testing) SelectFromFilter(s Filter) (interface{}, error) {
	ret := []Testing{}
	return ret, s.Select(&ret)
}

func (Testing) CountFromFilter(s Filter) (int64, error) {
	return s.Count(Testing{})
}

func (Testing) FromIdSave(data map[string]interface{}) (error, interface{}) {
	t := Testing{}
	err := mapstructure.WeakDecode(data, &t)
	if err != nil {
		return fmt.Errorf("FromIdSave WeakDecode error %s", err.Error()), t
	}
	if t.EnvId <= 0 {
		return fmt.Errorf("EnvId empty"), t
	}
	if t.Name == "" {
		return fmt.Errorf("Name empty"), t
	}
	if t.Project == "" {
		return fmt.Errorf("Project empty"), t
	}
	if t.Data == "" {
		return fmt.Errorf("Data empty"), t
	}

	if t.ID <= 0 {
		err = GormDB().Create(&t).Error
		if err != nil {
			return fmt.Errorf("Testing Create Error %s ", err.Error()), t
		}
		return nil, t
	}

	existT := Testing{}
	exist, err := GormFirst(&existT, &Testing{ID: t.ID}, "ID")
	if err != nil {
		return fmt.Errorf("FromIdSave GormFirst %s", err.Error()), t
	}
	if !exist {
		return fmt.Errorf("Testing %d not exist", t.ID), t
	}
	err = GormDB().Model(&t).Select("EnvId", "UserId", "Name", "Project", "Data").Updates(t).Error
	if err != nil {
		return fmt.Errorf("Testing Updates error %s ", err.Error()), t
	}
	return nil, t
}

func (Testing) FromIdDel(id uint) error {
	err := GormDB().Delete(&Testing{ID: id}).Error
	if err != nil {
		return fmt.Errorf("Testing FromIdDel Delete error %s ", err.Error())
	}
	return nil
}

func (Testing) FromApiIdDel(id uint) error {
	err := GormDB().Where(&Testing{ApiId: id}, "ApiId").Delete(Testing{}).Error
	if err != nil {
		return fmt.Errorf("Testing FromApiIdDel Delete error %s ", err.Error())
	}
	return nil
}
