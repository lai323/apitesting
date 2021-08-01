package entity

import (
	"fmt"

	"gorm.io/gorm/schema"
)

type TableDataIniter interface {
	schema.Tabler
	DataInit() error
}

var ToCreateTables = []TableDataIniter{
	&User{},
	&Env{},
	&Api{},
	&Testing{},
}

func Install() error {
	for _, t := range ToCreateTables {
		var err error
		if GormDB().Migrator().HasTable(t) {
			return fmt.Errorf("table %s already exist", t.TableName())
		}
		err = GormDB().Migrator().CreateTable(t)
		if err != nil {
			panic(fmt.Errorf("create table %s err :%s", t.TableName(), err.Error()))
		}
		err = t.DataInit()
		if err != nil {
			panic(fmt.Errorf("init table data %s err :%s", t.TableName(), err.Error()))
		}
	}

	return nil
}
