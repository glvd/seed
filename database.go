package seed

import (
	"github.com/go-xorm/xorm"
)

// Database ...
type Database struct {
	eng       *xorm.Engine
	syncTable []interface{}
	writer    chan SQLWriter
}

var _ Optioner = &Database{}

// NewDatabase ...
func NewDatabase(eng *xorm.Engine) *Database {
	db := new(Database)
	db.eng = eng
	db.writer = make(chan SQLWriter, 10)
	return db
}

// PushWriter ...
func (db *Database) PushWriter(s *xorm.Session, v interface{}) {
	db.writer <- sqlWriter(s, v)
}

// Sync ...
func (db *Database) Sync() error {
	if db.syncTable == nil {
		return nil
	}
	return db.eng.Sync2(db.syncTable...)
}

// RegisterSync ...
func (db *Database) RegisterSync(v interface{}) {
	db.syncTable = append(db.syncTable, v)
}

// Option ...
func (db *Database) Option(seed *Seed) {
	databaseOption(db)(seed)
}

// WriteCallback ...
type WriteCallback func(able SQLUpdateAble) (id interface{}, e error)

type write struct {
	cb      WriteCallback
	update  bool
	session *xorm.Session
	data    interface{}
	able    SQLUpdateAble
}

// UnfinishedWriter ...
func sqlWriter(session *xorm.Session, v interface{}) SQLWriter {
	return &write{
		session: session,
		data:    v,
	}
}

// Insert ...
func (w *write) Insert() (int64, error) {
	if w.update {
		_, e := w.cb(w.able)
		if e != nil {
			return 0, e
		}
		return w.session.ID(w.able.GetID()).Update(w.able)
	}
	return w.session.Insert(w.data)
}

// DatabaseArgs ...
type DatabaseArgs func(*Database)

// DatabaseShowSQLArg ...
func DatabaseShowSQLArg() DatabaseArgs {
	return func(db *Database) {
		db.eng.ShowSQL()
	}
}

// DatabaseShowExecTimeArg ...
func DatabaseShowExecTimeArg() DatabaseArgs {
	return func(db *Database) {
		db.eng.ShowExecTime()
	}
}
