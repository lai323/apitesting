package entity

import (
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/spf13/viper"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var gormdb *gorm.DB
var gormdbOnce sync.Once

func GormDB() *gorm.DB {

	gormdbOnce.Do(func() {
		var err error
		gormdb, err = CreateGormDB(
			viper.GetString("sqlite.dbfile"),
			viper.GetString("sqlite.loglevel"),
		)
		if err != nil {
			panic(err)
		}
	})
	return gormdb
}

func CreateGormDB(dbfile, loglevel string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dbfile), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      ConvDBLogLevel(loglevel),
			Colorful:      true,
		}),
	})
}

var loglevel = map[string]logger.LogLevel{
	"silent": logger.Silent,
	"error":  logger.Error,
	"warn":   logger.Warn,
	"info":   logger.Info,
}

func ConvDBLogLevel(level string) logger.LogLevel {
	if l, ok := loglevel[level]; ok {
		return l
	}
	return logger.Warn
}

func GormFirst(record, query interface{}, args ...interface{}) (bool, error) {

	err := GormDB().Where(query, args...).First(record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
