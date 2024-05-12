package plan

import (
	"context"
	"fmt"
	"time"
)

type TimeoutExecutableWrapper struct {
	Timeout  time.Duration
	Delegate Executable
}

func (exec TimeoutExecutableWrapper) Exec(ctx context.Context, opts ExecOption) error {
	if exec.Timeout <= 0 {
		return exec.Delegate.Exec(ctx, opts)
	}
	timoutCtx, cancelFn := context.WithTimeout(ctx, exec.Timeout)
	defer cancelFn()
	return exec.Delegate.Exec(timoutCtx, opts)
}

func (exec TimeoutExecutableWrapper) String() string {
	return fmt.Sprintf(`%v`, exec.Delegate)
}
