package seed

import (
	"github.com/glvd/seed/model"
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
func (db *Database) PushWriter(w SQLWriter) {
	db.writer <- w
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

type unfinishedWriter struct {
	session    *xorm.Session
	unfinished *model.Unfinished
}

// UnfinishedWriter ...
func UnfinishedWriter(session *xorm.Session, u *model.Unfinished) SQLWriter {
	return &unfinishedWriter{
		session:    session,
		unfinished: u,
	}
}

// Insert ...
func (w *unfinishedWriter) Insert() (int64, error) {
	return w.session.Insert(w.unfinished)
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
