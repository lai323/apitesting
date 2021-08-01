package entity

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type FilterSelecter interface {
	schema.Tabler
	CountFromFilter(Filter) (int64, error)
	SelectFromFilter(Filter) (interface{}, error)
}

type FromIdSaver interface {
	schema.Tabler
	FromIdSave(map[string]interface{}) (error, interface{})
}

type FromIdDeleter interface {
	schema.Tabler
	FromIdDel(id uint) error
}

type Operate string

var (
	OPERATE_EQ   Operate = "="
	OPERATE_NE   Operate = "<>"
	OPERATE_GT   Operate = ">"
	OPERATE_GTE  Operate = ">="
	OPERATE_LT   Operate = "<"
	OPERATE_LTE  Operate = "<="
	OPERATE_LIKE Operate = "LIKE"
)

type Filter struct {
	Sort   string
	Order  string
	Offset int
	Limit  int
	Where  map[string]string
	Op     map[string]Operate
}

func (s Filter) txWithWhere(tx *gorm.DB) (*gorm.DB, error) {
	for field, value := range s.Where {
		op, ok := s.Op[field]
		if !ok {
			op = OPERATE_EQ
		}
		switch op {
		case OPERATE_EQ:
			fallthrough
		case OPERATE_NE:
			fallthrough
		case OPERATE_GT:
			fallthrough
		case OPERATE_GTE:
			fallthrough
		case OPERATE_LT:
			fallthrough
		case OPERATE_LTE:
			tx = tx.Where(fmt.Sprintf("`%s` %s ?", field, op), value)
		case OPERATE_LIKE:
			tx = tx.Where(fmt.Sprintf("`%s` %s ?", field, op), fmt.Sprintf("%%%s%%", value))
		default:
			return tx, fmt.Errorf("Select not support operate %s", op)
		}
	}
	return tx, nil
}

func (s Filter) txWithSort(tx *gorm.DB) *gorm.DB {
	if s.Sort != "" {
		tx = tx.Order(fmt.Sprintf("`%s` %s", s.Sort, s.Order))
	}
	return tx
}

func (s Filter) txWithLimit(tx *gorm.DB) *gorm.DB {
	if s.Limit != 0 {
		tx = tx.Limit(s.Limit)
	}
	return tx
}

func (s Filter) txWithOffset(tx *gorm.DB) *gorm.DB {
	if s.Offset != 0 {
		tx = tx.Offset(s.Offset)
	}
	return tx
}

func (s Filter) getTx() (*gorm.DB, error) {
	tx := GormDB()
	var err error
	tx, err = s.txWithWhere(tx)
	if err != nil {
		return tx, err
	}
	if tx.Error != nil {
		return tx, tx.Error
	}

	tx = s.txWithSort(tx)
	if tx.Error != nil {
		return tx, tx.Error
	}
	tx = s.txWithLimit(tx)
	if tx.Error != nil {
		return tx, tx.Error
	}
	tx = s.txWithOffset(tx)
	if tx.Error != nil {
		return tx, tx.Error
	}
	return tx, tx.Error
}

func (s Filter) getColumnName(field string, structType reflect.Type) (string, error) {
	f, fexist := structType.FieldByName(field)
	if !fexist {
		return "", fmt.Errorf("fastadmin Select getColumnName %s can not find field %s", structType.Name(), field)
	}

	columnName, hasColumnName := schema.ParseTagSetting(f.Tag.Get("gorm"), ";")["COLUMN"]
	if !hasColumnName {
		columnName = schema.NamingStrategy{}.ColumnName("", field)
	}
	return columnName, nil
}

// 将 Filter, Op, Sort 中的 struct 字段名，转化为数据库字段名
func (s *Filter) convAllColumnName(structType reflect.Type) error {
	newFilter := map[string]string{}
	for fname, fvalue := range s.Where {
		columnName, err := s.getColumnName(fname, structType)
		if err != nil {
			return err
		}
		newFilter[columnName] = fvalue
	}
	s.Where = newFilter

	newOp := map[string]Operate{}
	for fname, fvalue := range s.Op {
		columnName, err := s.getColumnName(fname, structType)
		if err != nil {
			return err
		}
		newOp[columnName] = fvalue
	}
	s.Op = newOp

	if s.Sort != "" {
		var err error
		s.Sort, err = s.getColumnName(s.Sort, structType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Filter) Select(instances interface{}) error {
	if reflect.ValueOf(instances).Elem().Kind() != reflect.Slice {
		return fmt.Errorf("fastadmin Select Select instances must be pointer of model struct slice ")
	}
	err := s.convAllColumnName(reflect.TypeOf(instances).Elem().Elem())
	if err != nil {
		return err
	}

	tx, err := s.getTx()
	if err != nil {
		return err
	}
	tx.Find(instances)
	return tx.Error
}

func (s Filter) Count(t interface{}) (int64, error) {
	var count int64
	err := s.convAllColumnName(reflect.TypeOf(t))
	if err != nil {
		return count, err
	}

	tx := GormDB()
	tx, err = s.txWithWhere(tx)
	if err != nil {
		return count, err
	}
	if tx.Error != nil {
		return count, tx.Error
	}

	err = tx.Model(t).Count(&count).Error
	return count, err
}
