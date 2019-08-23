package seed

import (
	"context"
	"time"

	"github.com/glvd/seed/model"
	"github.com/go-xorm/xorm"
	"go.uber.org/atomic"
	"golang.org/x/xerrors"
)

// Database ...
type Database struct {
	Seeder
	eng       *xorm.Engine
	syncTable []interface{}
	cb        chan DatabaseCaller
	done      chan bool
	state     *atomic.Int32
}

// State ...
func (db *Database) State() State {
	return State(db.state.Load())
}

// Done ...
func (db *Database) Done() <-chan bool {
	go func() {
		db.cb <- nil
	}()
	return db.done
}

var _ DatabaseCaller = &databaseCall{}

type databaseCall struct {
	v  interface{}
	cb DatabaseCallbackFunc
}

// Call ...
func (c *databaseCall) Call(database *Database, eng *xorm.Engine) (e error) {
	return c.cb(database, eng, c.v)
}

// DatabaseCallback ...
func DatabaseCallback(v interface{}, cb DatabaseCallbackFunc) DatabaseCaller {
	return &databaseCall{
		v:  v,
		cb: cb,
	}
}

// Push ...
func (db *Database) Push(v interface{}) error {
	return db.pushDatabaseCallback(v)
}

// Run ...
func (db *Database) Run(ctx context.Context) {
	log.Info("database running")
	var e error
	e = db.Sync()
	if e != nil {
		panic(e)
	}
DatabaseEnd:
	for {
		select {
		case <-ctx.Done():
			log.Info("context end")
			db.state.Store(int32(StateStop))
			return
		case v := <-db.cb:
			db.state.Store(int32(StateRunning))
			if v == nil {
				log.Info("db end")
				break DatabaseEnd
			}
			e = v.Call(db, db.eng)
			if e != nil {
				log.Error(e)
			}
		case <-time.After(30 * time.Second):
			db.state.Store(int32(StateWaiting))
		}
	}
	close(db.cb)
	db.done <- true
}

// BeforeRun ...
func (db *Database) BeforeRun(seed Seeder) {
	db.Seeder = seed
}

// AfterRun ...
func (db *Database) AfterRun(seed Seeder) {
}

// NewDatabase ...
func NewDatabase(eng *xorm.Engine, args ...DatabaseArgs) *Database {
	db := new(Database)
	db.eng = eng
	db.cb = make(chan DatabaseCaller, 10)
	db.done = make(chan bool)
	for _, argFn := range args {
		argFn(db)
	}

	return db
}

// PushCallback ...
func (db *Database) pushDatabaseCallback(cb interface{}) (e error) {
	if v, b := cb.(DatabaseCaller); b {
		db.cb <- v
		return nil
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
func (db *Database) RegisterSync(v ...interface{}) {
	for _, val := range v {
		db.syncTable = append(db.syncTable, val)
	}
}

// Option ...
func (db *Database) Option(seed Seeder) {
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

// databaseOption ...
func databaseOption(db *Database) Options {
	return func(seed Seeder) {
		seed.SetThread(StepperDatabase, db)
		//seed.SetNormalThread(StepperRDatabase, db)
	}
}
