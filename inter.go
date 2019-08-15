package seed

import "github.com/go-xorm/xorm"

// Seeder ...
type Seeder interface {
	Start()
	Wait()
	Stop()
	Err() error
}

// SQLWriter ...
type SQLWriter interface {
	Insert() (int64, error)
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
