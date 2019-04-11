package model

import (
	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
	"reflect"
	"time"
)

var db *xorm.Engine
var syncTable map[string]interface{}

// RegisterTable ...
func RegisterTable(v interface{}) {
	tof := reflect.TypeOf(v).Name()
	if syncTable == nil {
		syncTable = map[string]interface{}{
			tof: v,
		}
	}
	syncTable[tof] = v
}

// DB ...
func DB() *xorm.Engine {
	if db == nil {
		if err := InitDB(); err != nil {
			panic(err)
		}
	}
	return db
}

// InitDB ...
func InitDB() (e error) {
	eng, e := xorm.NewEngine("sqlite3", "seed.db")
	if e != nil {
		return e
	}
	eng.ShowSQL(true)
	eng.ShowExecTime(true)

	for _, val := range syncTable {
		e := eng.Sync2(val)
		if e != nil {
			return e
		}
	}

	db = eng
	return nil
}

// Model ...
type Model struct {
	ID        string     `json:"-" xorm:"id pk"`
	CreatedAt time.Time  `json:"-" xorm:"created_at created"`
	UpdatedAt time.Time  `json:"-" xorm:"updated_at updated"`
	DeletedAt *time.Time `json:"-" xorm:"deleted_at deleted"`
	//Version   int        `json:"-" xorm:"version"`
}

// BeforeInsert ...
func (m *Model) BeforeInsert() {
	if m.ID == "" {
		m.ID = uuid.Must(uuid.NewRandom()).String()
	}
}
