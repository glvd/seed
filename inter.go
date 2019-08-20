package seed

import "github.com/go-xorm/xorm"

// Seeder ...
type Seeder interface {
	Start()
	Wait()
	Stop()
	Err() error
}

// SQLUpdateAble ...
type SQLUpdateAble interface {
	GetID() string
	SetID(string)
	GetVersion() int
	SetVersion(int)
}

// DatabaseCallback ...
type DatabaseCallback func(database *Database, eng *xorm.Engine) (e error)

// SQLWriter ...
type SQLWriter interface {
	InsertOrUpdate() (int64, error)
	Done()
	Failed()
}

// SQLReader ...
type SQLReader interface {
	FindOne(*xorm.Session, interface{}) error
	FindAll(*xorm.Session, interface{}) error
}

// Initer ...
type Initer interface {
	Init()
}

//Optioner set option
type Optioner interface {
	Option(seed *Seed)
}
