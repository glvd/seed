package seed

import "github.com/go-xorm/xorm"

// Database ...
type Database struct {
	eng *xorm.Engine
}

// NewDatabase ...
func NewDatabase(eng *xorm.Engine) *Database {
	db := new(Database)
	db.eng = eng
	return db
}
