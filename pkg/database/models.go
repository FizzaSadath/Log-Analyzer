package database

import (
	"time"

	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model
	TimeStamp time.Time
	Level     string
	Component string
	Host      string
	RequestId string
	Message   string
}
