package tabler

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type table struct {
	db        *gorm.DB
	cursorKey string
	Pager
}

type Pager struct {
	PageSize int32
	PageNum  int32
	Cursor   int64
}

func NewTable(db *gorm.DB, cursorKey string, pager Pager) *table {
	return &table{
		db:        db,
		cursorKey: cursorKey,
		Pager:     pager,
	}
}

//Build 会比 t.PageSize 多查出一条数据，用来判断是否还有下一页数据
//调用方需要做此判断
func (t table) Build() (res *gorm.DB) {
	res = t.db
	if t.PageNum == 0 {
		res = t.db.Where(fmt.Sprintf("%s < ?", t.cursorKey), t.Cursor)
		res = res.Limit(t.PageSize)
		return res
	}
	offset := t.PageSize * (t.PageNum - 1)
	//t.PageSize + 1 多查出一条用来判断是否还有下一页数据
	//调用方需要做此判断
	return res.Offset(offset).Limit(t.PageSize + 1)
}

func EmptyDBFuc(db *gorm.DB) *gorm.DB {
	return db
}

//Do 会比 t.PageSize 多查出一条数据，用来判断是否还有下一页数据
//调用方需要做此判断
func (t table) Do(doer func(db *gorm.DB) *gorm.DB) error {
	return doer(t.Build()).Error
}
