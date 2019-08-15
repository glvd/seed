package seed

import "github.com/go-xorm/xorm"

// Database ...
type Database struct {
	eng       *xorm.Engine
	syncTable []interface{}
}

// NewDatabase ...
func NewDatabase(eng *xorm.Engine) *Database {
	db := new(Database)
	db.eng = eng
	return db
}

// Sync ...
func (db *Database) Sync() error {
	if db.syncTable == nil {
		return nil
	}
	return db.eng.Sync2(db.syncTable...)
}
