package awscli

import (
	"context"
	"time"
)

type ValidatorAdapter struct {
	Executor CommandExecutor
}

func (v *ValidatorAdapter) Validate(profileName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := v.Executor.Run(ctx, "aws", "sts", "get-caller-identity", "--profile", profileName)
	return err
}
