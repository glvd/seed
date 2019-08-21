package seed

import (
	"context"
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

// Threader ...
type Threader interface {
	Runnable
	Pusher
	BeforeRun(seed Seeder)
	AfterRun(seed Seeder)
}

// Runnable ...
type Runnable interface {
	Run(context.Context)
}

// Pusher ...
type Pusher interface {
	Push(interface{}) error
}
