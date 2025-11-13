package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (l Entry) String() string {
	if l.TimeStamp.IsZero() {
		return "Empty"
	} else {
		return fmt.Sprintf("%s : %s : %s : %s : %s : %s", l.TimeStamp, l.Level, l.Component, l.Host, l.RequestId, l.Message)

	}
}

func CreateDB(dbUrl string) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level (Silent, Error, Warn, Info)
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error
			Colorful:                  true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, fmt.Errorf("couldn't open database %v", err)
	}
	return db, nil

}
func InitDB(db *gorm.DB) error {
	db.AutoMigrate(&Entry{})
	return nil
}

func AddDB(db *gorm.DB, e Entry) error {
	ctx := context.Background()
	err := gorm.G[Entry](db).Create(ctx, &e)
	if err != nil {
		return err
	}
	return nil
}

// func QueryDB(db *gorm.DB, query string) ([]Entry, error) {
// 	//parse the query
// 	var ret []Entry
// 	parsed
// }
