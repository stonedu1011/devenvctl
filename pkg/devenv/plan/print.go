package plan

import (
	"context"
)

type PrintExecutable string

func (exec PrintExecutable) Exec(ctx context.Context) error {
	logger.WithContext(ctx).Infof(string(exec))
	return nil
}

func (exec PrintExecutable) String() string {
	return "print: " + string(exec)
}
