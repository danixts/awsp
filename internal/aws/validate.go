package aws

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func ValidateProfile(profile string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "aws", "sts", "get-caller-identity", "--profile", profile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf("%w: %s", err, stderr.String())
		}
		return err
	}
	if stdout.Len() > 0 {
		_, _ = os.Stderr.Write(stdout.Bytes())
	}
	return nil
}
