package seed

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/xerrors"
)

// Move ...
type Move struct {
	*Thread
	cb     chan MoveCaller
	ToPath string
}

// Register ...
func (m *Move) Register(ops ...Optioner) {

}

// Push ...
func (m *Move) Push(v interface{}) error {
	return m.push(v)
}

func (m *Move) push(v interface{}) error {
	if mc, b := v.(MoveCaller); b {
		m.cb <- mc
		return nil
	}
	return xerrors.New("not move callback")
}

// NewMove ...
func NewMove() *Move {
	return &Move{
		Thread: NewThread(),
	}
}

// Run ...
func (m *Move) Run(ctx context.Context) {
	log.Info("move running")

InfoEnd:
	for {
		select {
		case <-ctx.Done():
			break InfoEnd
		case cb := <-m.cb:
			if cb == nil {
				break InfoEnd
			}
			m.SetState(StateRunning)
			e := cb.Call(m)
			if e != nil {
				log.Error(e)
			}
		case <-time.After(30 * time.Second):
			log.Info("move time out")
			m.SetState(StateWaiting)
		}
	}
	close(m.cb)
	m.Finished()
}

// MoveOption ...
func MoveOption(stepper Stepper, Move *Move) Options {
	return func(seed Seeder) {
		seed.SetBaseThread(stepper, Move)
	}
}

type moveCall struct {
	fromPath string
	toPath   string
}

// Call ...
func (m *moveCall) Call(*Move) error {
	inputFile, err := os.Open(m.fromPath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	defer inputFile.Close()
	outputFile, err := os.Create(m.toPath)
	if err != nil {
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(m.fromPath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}
	return nil
}

// MoveCall ...
func MoveCall(from, to string) (Stepper, MoveCaller) {
	return StepperMove, &moveCall{fromPath: from, toPath: to}

}
