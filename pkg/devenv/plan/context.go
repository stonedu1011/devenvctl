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

var DefaultExecOption ExecOption

const (
	ActionStart   Action = "start"
	ActionStop    Action = "stop"
	ActionRestart Action = "restart"
)

type Action string

type ExecutionPlanner interface {
	Plan(action Action) (ExecutionPlan, error)
}

type ExecOptions func(opt *ExecOption)
type ExecOption struct {
	Verbose bool
	DryRun  bool
}

type Executable interface {
	Exec(ctx context.Context, opts ExecOption) error
}

type ExecutionPlan interface {
	Steps() []Executable
	Metadata() interface{}
	Execute(ctx context.Context, opts...ExecOptions) error
}

func NewExecutionPlan(metadata interface{}, execs ...Executable) ExecutionPlan {
	return execPlan{
		steps: execs,
		meta:  metadata,
	}
}

type execPlan struct {
	steps []Executable
	meta  interface{}
}

func (p execPlan) Steps() []Executable {
	return p.steps
}

func (p execPlan) Metadata() interface{} {
	return p.meta
}

func (p execPlan) Execute(ctx context.Context, opts...ExecOptions) error {
	opt := DefaultExecOption
	for _, fn := range opts {
		fn(&opt)
	}

	if opt.DryRun {
		return p.DryRun(ctx)
	}
	for _, exec := range p.steps {
		if e := exec.Exec(ctx, opt); e != nil {
			return e
		}
	}
	return nil
}

func (p execPlan) DryRun(ctx context.Context) error {
	if len(p.steps) == 0 {
		logger.WithContext(ctx).Infof("DryRun - Planned Steps: NONE")
		return nil
	}
	logger.WithContext(ctx).Infof("DryRun - planned steps:")
	for _, exec := range p.steps {
		fmt.Printf(`    %v\n`, exec)
	}
	return nil
}

func NewClosableExecutionPlan(metadata interface{}, closerFunc func() error, execs ...Executable) ExecutionPlan {
	return closerPlan{
		ExecutionPlan: NewExecutionPlan(metadata, execs...),
		closeFunc: closerFunc,
	}
}

type closerPlan struct {
	ExecutionPlan
	closeFunc func() error
}

func (p closerPlan) Close() error {
	if p.closeFunc != nil {
		return p.closeFunc()
	}
	return nil
}
