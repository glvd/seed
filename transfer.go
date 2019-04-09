package seed

import (
	"github.com/jinzhu/gorm"
)

var lite = DB()

// DB ...
func DB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "sedd.db")
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

// Transfer ...
func Transfer() {
	//for _, val := range VideoList {
	//
	//}
}
