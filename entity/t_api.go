package entity

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Api struct {
	ID        uint   `gorm:"primaryKey"`
	Project   string `gorm:"not null;comment:项目"`
	Name      string `gorm:"not null;comment:名称"`
	Route     string `gorm:"not null;comment:路由"`
	Define    string `gorm:"type:text;not null;comment:环境变量"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Api) TableName() string {
	return "apitesting_api"
}

func (Api) DataInit() error {
	return nil
}

func (Api) SelectFromFilter(s Filter) (interface{}, error) {
	ret := []Api{}
	return ret, s.Select(&ret)
}

func (Api) CountFromFilter(s Filter) (int64, error) {
	return s.Count(Api{})
}

func (Api) FromIdSave(data map[string]interface{}) (error, interface{}) {
	api := Api{}
	err := mapstructure.WeakDecode(data, &api)
	if err != nil {
		return fmt.Errorf("FromIdSave WeakDecode error %s", err.Error()), api
	}
	if api.Project == "" {
		return fmt.Errorf("Project empty"), api
	}
	if api.Name == "" {
		return fmt.Errorf("Name empty"), api
	}
	if api.Route == "" {
		return fmt.Errorf("Route empty"), api
	}
	if api.Define == "" {
		return fmt.Errorf("Define empty"), api
	}

	if api.ID <= 0 {
		existApi := Api{}
		exist, err := GormFirst(&existApi, &Api{Project: api.Project, Route: api.Route}, "Project", "Route")
		if err != nil {
			return fmt.Errorf("FromIdSave GormFirst %s", err.Error()), api
		}
		if exist {
			return fmt.Errorf("Api %s %s exist", api.Project, api.Route), api
		}

		err = GormDB().Create(&api).Error
		if err != nil {
			return fmt.Errorf("Api Create Error %s ", err.Error()), api
		}
		return nil, api
	}

	existApi := Api{}
	exist, err := GormFirst(&existApi, &Api{ID: api.ID}, "ID")
	if err != nil {
		return fmt.Errorf("FromIdSave GormFirst %s", err.Error()), api
	}
	if !exist {
		return fmt.Errorf("Api %d not exist", api.ID), api
	}
	err = GormDB().Model(&api).Select("Project", "Name", "Route", "Define").Updates(api).Error
	if err != nil {
		return fmt.Errorf("Api Updates error %s ", err.Error()), api
	}
	return nil, api
}

func (Api) FromIdDel(id uint) error {
	err := GormDB().Delete(&Api{ID: id}).Error
	if err != nil {
		return fmt.Errorf("Api Delete error %s ", err.Error())
	}
	return Testing{}.FromApiIdDel(id)
}
