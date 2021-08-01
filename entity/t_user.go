package entity

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique;not null;comment:用户名"`
	Password  string `gorm:"not null;comment:密码"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "apitesting_user"
}

// 创建系统账户
func (User) DataInit() error {
	users := []User{}
	userscfg := viper.Get("apitesting.default_account").([]interface{})
	err := mapstructure.WeakDecode(userscfg, &users)
	if err != nil {
		return err
	}
	return GormDB().Create(&users).Error
}

func (User) SelectFromFilter(s Filter) (interface{}, error) {
	ret := []User{}
	return ret, s.Select(&ret)
}

func (User) CountFromFilter(s Filter) (int64, error) {
	return s.Count(User{})
}

func (User) FromIdSave(data map[string]interface{}) (error, interface{}) {
	u := User{}
	err := mapstructure.WeakDecode(data, &u)
	if err != nil {
		return fmt.Errorf("FromIdSave WeakDecode error %s", err.Error()), u
	}
	if u.Password == "" {
		return fmt.Errorf("Password empty"), u
	}

	if u.ID <= 0 {
		err = GormDB().Create(&u).Error
		if err != nil {
			return fmt.Errorf("User Create Error %s ", err.Error()), u
		}
		return nil, u
	}

	existT := User{}
	exist, err := GormFirst(&existT, &User{ID: u.ID}, "ID")
	if err != nil {
		return fmt.Errorf("FromIdSave GormFirst %s", err.Error()), u
	}
	if !exist {
		return fmt.Errorf("User %d not exist", u.ID), u
	}
	err = GormDB().Model(&u).Select("Password").Updates(u).Error
	if err != nil {
		return fmt.Errorf("User Updates error %s ", err.Error()), u
	}
	return nil, u
}

func (User) FromIdDel(id uint) error {
	adminid, err := User{}.AdminId()
	if err != nil {
		return err
	}
	if id == adminid {
		return fmt.Errorf("Can't Delete admin")
	}
	err = GormDB().Delete(&User{ID: id}).Error
	if err != nil {
		return fmt.Errorf("User Delete error %s ", err.Error())
	}
	return nil
}

func (User) AdminId() (uint, error) {
	u := User{}
	exist, err := GormFirst(&u, &User{Name: "admin"}, "Name")
	if err != nil {
		return 0, fmt.Errorf("AdminId GormFirst %s", err.Error())
	}
	if !exist {
		return 0, fmt.Errorf("AdminId not found admin")
	}
	return u.ID, nil
}
