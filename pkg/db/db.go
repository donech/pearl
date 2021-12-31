//Package db code
package db

import (
	"context"
	"fmt"
	"time"

	logger "github.com/donech/pearl/pkg/log"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var CreatedFiledName = "created_time"
var UpdatedFiledName = "updated_time"
var DeletedFiledName = "deleted_time"

//New just like Open
func New(cfg Config) (*gorm.DB, func()) {
	return Open(cfg)
}

//Open the database connection
func Open(conf Config) (*gorm.DB, func()) {
	if db, err := gorm.Open("mysql", conf.Dsn); err != nil {
		panic(errors.WithStack(err))
	} else {
		db.LogMode(conf.LogMode)
		db.DB().SetMaxIdleConns(conf.MaxIdle)
		db.DB().SetMaxOpenConns(conf.MaxOpen)
		db.DB().SetConnMaxLifetime(time.Duration(conf.MaxLifetime) * time.Second)
		RegisterCallback(db)
		cleanup := func() {
			if err := db.Close(); err != nil {
				panic(errors.WithStack(err))
			}
		}

		return db, cleanup
	}
}

// Trace return clone DB with trace ctx
func Trace(ctx context.Context, db *gorm.DB) *gorm.DB {
	return db.Set("ctx", ctx)
}

func RegisterCallback(db *gorm.DB) {
	db.Callback().Create().Register("logger", func(scope *gorm.Scope) {
		log("create", scope)
	})
	db.Callback().Update().Register("logger", func(scope *gorm.Scope) {
		log("update", scope)
	})
	db.Callback().Delete().Register("logger", func(scope *gorm.Scope) {
		log("delete", scope)
	})
	db.Callback().Query().Register("logger", func(scope *gorm.Scope) {
		log("query", scope)
	})
	db.Callback().RowQuery().After("gorm:row_query").Register("logger", func(scope *gorm.Scope) {
		log("raw_query", scope)
	})
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
}

func log(msg string, scope *gorm.Scope) {
	ctx, ok := scope.Get("ctx")
	if !ok {
		return
	}
	c, ok := ctx.(context.Context)
	if !ok {
		return
	}
	logger.L(c).Info(
		msg,
		zap.String("sql", scope.SQL),
		zap.Reflect("vars", scope.SQLVars),
		zap.Reflect("result", scope.Value),
	)
}

// updateTimeStampForCreateCallback will set `CreatedTime`, `UpdatedTime` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now()
		if createTimeField, ok := scope.FieldByName(CreatedFiledName); ok {
			if createTimeField.IsBlank {
				err := createTimeField.Set(nowTime)
				if err != nil {
					log("set "+CreatedFiledName+" failed: "+err.Error(), scope)
				}
			}
		}

		if updatedField, ok := scope.FieldByName(UpdatedFiledName); ok {
			if updatedField.IsBlank {
				err := updatedField.Set(nowTime)
				if err != nil {
					log("set "+UpdatedFiledName+" failed: "+err.Error(), scope)
				}
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		err := scope.SetColumn(UpdatedFiledName, time.Now())
		if err != nil {
			log("set "+UpdatedFiledName+" failed: "+err.Error(), scope)
		}
	}
}

// deleteCallback will set `DeletedOn` where deleting
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedField, hasDeletedField := scope.FieldByName(DeletedFiledName)

		if !scope.Search.Unscoped && hasDeletedField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedField.DBName),
				scope.AddToVars(time.Now()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

// addExtraSpaceIfExist adds a separator
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
