package db

import (
	"database/sql/driver"
	"fmt"
	"time"
)

var DateTimeFormat = "2006-01-02 15:04:05"

type Entity struct {
	ID int64 `gorm:"primary_key" json:"id"`
}

type CUDTime struct {
	CreatedTime DateTime `json:"created_time"`
	UpdatedTime DateTime `json:"updated_time"`
	DeletedTime DateTime `json:"deleted_time"`
}

type DateTime struct {
	time.Time
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("\"0000-00-00 00:00:00\""), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format(DateTimeFormat))
	return []byte(formatted), nil
}

func (t DateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *DateTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = DateTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
