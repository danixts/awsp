package awscli

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type CommandExecutor interface {
	Run(ctx context.Context, name string, args ...string) ([]byte, error)
}

type CLIExecutor struct{}

func (e *CLIExecutor) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return nil, fmt.Errorf("%s: %s", name, stderr.String())
		}
		return nil, err
	}
	return stdout.Bytes(), nil
}

func defaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}
