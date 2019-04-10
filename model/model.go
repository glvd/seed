package model

import (
	"github.com/go-xorm/xorm"
	"time"
)

var database *xorm.Engine

// InitDB ...
func init() {
	database = DB()
}

// DB ...
func DB() *xorm.Engine {
	db, err := xorm.NewEngine("sqlite3", "seed.db")
	if err != nil {
		panic(err)
	}

	e := db.Sync2(&Video{})
	if e != nil {
		return &xorm.Engine{}
	}

	return db
}

// Model ...
type Model struct {
	ID        string     `xorm:"id"`
	CreatedAt time.Time  `xorm:"created_at"`
	UpdatedAt time.Time  `xorm:"updated_at"`
	DeletedAt *time.Time `xorm:"deleted_at"`
	Version   int        `xorm:"version"`
}
