package seed

import (
	"context"

	"github.com/glvd/seed/model"
	"github.com/go-xorm/xorm"
	"golang.org/x/xerrors"
)

// Database ...
type Database struct {
	eng       *xorm.Engine
	syncTable []interface{}
	cb        chan DatabaseCallback
}

// Run ...
func (db *Database) Run(ctx context.Context) {
	var e error
	for {
		select {
		case <-ctx.Done():
			return
		case v := <-db.cb:
			e = v(db, db.eng)
			if e != nil {
				log.Error(e)
			}
		}
	}
}

// BeforeRun ...
func (db *Database) BeforeRun(seed *Seed) {
}

// AfterRun ...
func (db *Database) AfterRun(seed *Seed) {
}

var _ Optioner = &Database{}

// NewDatabase ...
func NewDatabase(eng *xorm.Engine, args ...DatabaseArgs) *Database {
	db := new(Database)
	db.eng = eng
	db.cb = make(chan DatabaseCallback, 10)

	for _, argFn := range args {
		argFn(db)
	}

	return db
}

// PushCallback ...
func (db *Database) PushCallback(cb interface{}) (e error) {
	if v, b := cb.(DatabaseCallback); b {
		go func(database *Database, databaseCallback DatabaseCallback) {
			database.cb <- databaseCallback
		}(db, v)
	}
	return xerrors.New("not database callback")
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

// BeforeUpdate ...
type BeforeUpdate func(session *xorm.Session, modeler model.Modeler) (id interface{}, e error)

type write struct {
	cb      BeforeUpdate
	update  bool
	session *xorm.Session
	model   model.Modeler
}

// UnfinishedWriter ...
func sqlWriter(session *xorm.Session, m model.Modeler) *write {
	return &write{
		session: session,
		model:   m,
	}
}

// Insert ...
func (w *write) InsertOrUpdate() (int64, error) {
	if w.update {
		id, e := w.cb(w.session, w.model)
		if e != nil {
			return 0, e
		}
		return w.session.ID(id).Update(w.model)
	}
	return w.session.Insert(w.model)
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
