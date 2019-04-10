package seed

import (
	"github.com/jinzhu/gorm"
)

var database = DB()

// DB ...
func DB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "seed.db")
	if err != nil {
		panic("failed to connect database")
	}

	for _, val := range VideoList {
		database.AutoMigrate(val)
	}

	return db
}

// Transfer ...
func Transfer() {

}
