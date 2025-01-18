package data

import "gorm.io/gorm"

type Storage struct {
	*gorm.DB
}
