package seed

import (
	"context"

	"github.com/go-xorm/xorm"
)

// Seeder ...
type Seeder interface {
	Start()
	Wait()
	Stop()
	PushTo(stepper Stepper, v interface{}) error
	Register(ops ...Optioner)
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
	Option(seed *seed)
}

// Stepper ...
type Stepper int

// StepperNone ...
const (
	// StepperNone ...
	StepperNone Stepper = iota
	//StepperAPI ...
	StepperAPI
	//StepperDatabase ...
	StepperDatabase
	// StepperInformation ...
	StepperInformation
	// StepperMoveInfo ...
	StepperMoveInfo
	// StepperProcess ...
	StepperProcess
	// StepperMoveproc ...
	StepperMoveproc
	// StepperTransfer ...
	StepperTransfer
	// StepperPin ...
	StepperPin
	// StepperCheck ...
	StepperCheck
	// StepperUpdate ...
	StepperUpdate
	// StepperMax ...
	StepperMax
)

// Threader ...
type Threader interface {
	Runnable
	Pusher
	BeforeRun(seed *seed)
	AfterRun(seed *seed)
}

// Runnable ...
type Runnable interface {
	Run(context.Context)
}

// Pusher ...
type Pusher interface {
	Push(interface{}) error
}
