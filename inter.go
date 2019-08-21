package seed

import (
	"context"

	"github.com/go-xorm/xorm"
	httpapi "github.com/ipfs/go-ipfs-http-client"
)

// Seeder ...
type Seeder interface {
	Start()
	Wait()
	Stop()
	PushTo(stepper Stepper, v interface{}) error
	GetThread(stepper Stepper) Threader
	SetThread(stepper Stepper, threader Threader)
	HasThread(stepper Stepper) bool
	SetBaseThread(stepper Stepper, threader Threader)
	IsBase(stepper Stepper) bool
	Register(ops ...Optioner)
	Err() error
}

// Initer ...
type Initer interface {
	Init()
}

//Optioner set option
type Optioner interface {
	Option(Seeder)
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

// DatabaseCallbackFunc ...
type DatabaseCallbackFunc func(database *Database, eng *xorm.Engine, v interface{}) (e error)

// DatabaseCaller ...
type DatabaseCaller interface {
	Call(database *Database, eng *xorm.Engine) (e error)
}

// APICallbackFunc ...
type APICallbackFunc func(api *API, api2 *httpapi.HttpApi, v interface{}) (e error)

// APICaller ...
type APICaller interface {
	Call(*API, *httpapi.HttpApi) error
}

// Threader ...
type Threader interface {
	Runnable
	Pusher
	BeforeRun(seed Seeder)
	AfterRun(seed Seeder)
}

// Async ...
type Async interface {
	NeedWait() bool
	IsRunning() bool
	Rerun()
	Stop()
}

// Runnable ...
type Runnable interface {
	Run(context.Context)
}

// Pusher ...
type Pusher interface {
	Push(interface{}) error
}
