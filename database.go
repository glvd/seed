package seed

import (
	"context"

	"github.com/glvd/seed/model"
	"github.com/go-xorm/xorm"
)

// Database ...
type Database struct {
	eng       *xorm.Engine
	syncTable []interface{}
	writer    chan SQLWriter
}

// Run ...
func (db *Database) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
		case v := <-db.writer:
			_, e := v.InsertOrUpdate()
			if e != nil {
				log.Error(e)
				v.Failed()
				continue
			}
			v.Done()
		}
	}
}

// BeforeRun ...
func (db *Database) BeforeRun(seed *Seed) {
	panic("implement me")
}

// AfterRun ...
func (db *Database) AfterRun(seed *Seed) {
	panic("implement me")
}

var _ Optioner = &Database{}

// NewDatabase ...
func NewDatabase(eng *xorm.Engine, args ...DatabaseArgs) *Database {
	db := new(Database)
	db.eng = eng
	db.writer = make(chan SQLWriter, 10)

	for _, argFn := range args {
		argFn(db)
	}

	return db
}

// PushWriter ...
func (db *Database) PushWriter(s *xorm.Session, v model.Modeler) {
	db.writer <- sqlWriter(s, v)
}

// PushCallbackWriter ...
func (db *Database) PushCallbackWriter(s *xorm.Session, v model.Modeler, callback BeforeUpdate) {
	w := sqlWriter(s, v)
	w.cb = callback
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

// BeforeUpdate ...
type BeforeUpdate func(session *xorm.Session, modeler model.Modeler) (id interface{}, e error)

type write struct {
	cb      BeforeUpdate
	update  bool
	session *xorm.Session
	model   model.Modeler
}

// Done ...
func (w *write) Done() {
	panic("implement me")
}

// Failed ...
func (w *write) Failed() {
	panic("implement me")
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
