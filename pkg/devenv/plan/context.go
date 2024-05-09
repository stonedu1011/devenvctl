package plan

import (
	"context"
	"errors"
	"fmt"
	"github.com/cisco-open/go-lanai/pkg/log"
)

var logger = log.New("CLI")

var (
	ErrPlanNotAvailable = errors.New(`plan for given action is not available`)
)

const (
	ActionStart   Action = "start"
	ActionStop    Action = "stop"
	ActionRestart Action = "restart"
)

type Action string

type ExecutionPlanner interface {
	Plan(action Action) ([]Executable, error)
}

type Executable interface {
	Exec(ctx context.Context) error
}

func Execute(ctx context.Context, executables ...Executable) error {
	for _, exec := range executables {
		if e := exec.Exec(ctx); e != nil {
			return e
		}
	}
	return nil
}

func DryRun(ctx context.Context, executables ...Executable) error {
	if len(executables) == 0 {
		logger.WithContext(ctx).Infof("DryRun - Planned Steps: NONE")
		return nil
	}
	logger.WithContext(ctx).Infof("DryRun - planned steps:")
	for _, exec := range executables {
		fmt.Println(exec)
	}
	return nil
}
