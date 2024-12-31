package database

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Date string

func (Date) GormDataType() string {
	return "date"
}

func (Date) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "date"
}

func (date Date) Value() (driver.Value, error) {
	return string(date), nil
}

func (date *Date) Scan(value interface{}) error {
	scanned, ok := value.(time.Time)
	if !ok {
		return errors.New(fmt.Sprintf("Failed to scan date:", value))
	}

	scannedString := scanned.Format("2006-01-02")
	*date = Date(scannedString)
	return nil
}
