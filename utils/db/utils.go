package db

import (
	"errors"
	"gorm.io/gorm"
)

func CheckErrIsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
