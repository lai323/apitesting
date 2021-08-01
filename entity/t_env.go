package entity

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Env struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique;not null;comment:环境名称"`
	Data      string `gorm:"type:text;not null;comment:环境变量数据"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Env) TableName() string {
	return "apitesting_env"
}

func (Env) DataInit() error {
	return nil
}

func (Env) SelectFromFilter(s Filter) (interface{}, error) {
	ret := []Env{}
	return ret, s.Select(&ret)
}

func (Env) CountFromFilter(s Filter) (int64, error) {
	return s.Count(Env{})
}

func (Env) FromIdSave(data map[string]interface{}) (error, interface{}) {
	env := Env{}
	err := mapstructure.WeakDecode(data, &env)
	if err != nil {
		return fmt.Errorf("FromIdSave WeakDecode error %s", err.Error()), env
	}
	if env.Name == "" {
		return fmt.Errorf("Name empty"), env
	}
	if env.Data == "" {
		return fmt.Errorf("Data empty"), env
	}

	if env.ID <= 0 {
		existEnv := Env{}
		exist, err := GormFirst(&existEnv, &Env{Name: env.Name}, "Name")
		if err != nil {
			return fmt.Errorf("FromIdSave GormFirst %s", err.Error()), env
		}
		if exist {
			return fmt.Errorf("FromIdSave env %s already exist", existEnv.Name), env
		}
		err = GormDB().Create(&env).Error
		if err != nil {
			return fmt.Errorf("Env Create Error %s ", err.Error()), env
		}
		return nil, env
	}

	existEnv := Env{}
	exist, err := GormFirst(&existEnv, &Env{ID: env.ID}, "ID")
	if err != nil {
		return fmt.Errorf("FromIdSave GormFirst %s", err.Error()), env
	}
	if !exist {
		return fmt.Errorf("Env %d not exist", env.ID), env
	}
	err = GormDB().Model(&env).Select("Name", "Data").Updates(env).Error
	if err != nil {
		return fmt.Errorf("Env Updates error %s ", err.Error()), env
	}
	return nil, env
}

func (Env) FromIdDel(id uint) error {
	return GormDB().Delete(&Env{ID: id}).Error
}
