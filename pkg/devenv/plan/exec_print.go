package plan

import (
	"context"
	"fmt"
)

type PrintExecutable string

func (exec PrintExecutable) Exec(ctx context.Context, opts ExecOption) error {
	if opts.DryRun {
		fmt.Printf("- %v\n", exec)
		return nil
	}
	logger.WithContext(ctx).Infof(string(exec))
	return nil
}

func (exec PrintExecutable) String() string {
	return "print: " + string(exec)
}
