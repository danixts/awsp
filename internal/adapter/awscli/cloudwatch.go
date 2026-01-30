package awscli

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/danixts/awsp/internal/domain"
)

type CloudWatchAdapter struct {
	Executor CommandExecutor
}

type describeLogGroupsResponse struct {
	LogGroups []domain.LogGroup `json:"logGroups"`
	NextToken *string           `json:"nextToken"`
}

func (a *CloudWatchAdapter) ListLogGroups(prefix string) ([]domain.LogGroup, error) {
	var all []domain.LogGroup
	var nextToken string

	for {
		args := []string{"logs", "describe-log-groups", "--output", "json"}
		if prefix != "" {
			args = append(args, "--log-group-name-prefix", prefix)
		}
		if nextToken != "" {
			args = append(args, "--next-token", nextToken)
		}

		ctx, cancel := defaultContext()
		out, err := a.Executor.Run(ctx, "aws", args...)
		cancel()
		if err != nil {
			return nil, fmt.Errorf("describe-log-groups: %w", err)
		}

		var resp describeLogGroupsResponse
		if err := json.Unmarshal(out, &resp); err != nil {
			return nil, fmt.Errorf("parse describe-log-groups: %w", err)
		}

		all = append(all, resp.LogGroups...)

		if resp.NextToken == nil || *resp.NextToken == "" {
			break
		}
		nextToken = *resp.NextToken
	}

	return all, nil
}

func (a *CloudWatchAdapter) TailLogs(ctx context.Context, logGroup, since, region string) (<-chan string, error) {
	ch := make(chan string, 100)
	go func() {
		defer close(ch)
		args := []string{"logs", "tail", logGroup, "--since", since, "--follow"}
		if region != "" {
			args = append(args, "--region", region)
		}
		cmd := exec.CommandContext(ctx, "aws", args...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			ch <- "Error: " + err.Error()
			return
		}
		cmd.Stderr = cmd.Stdout
		if err := cmd.Start(); err != nil {
			ch <- "Error: " + err.Error()
			return
		}
		sc := bufio.NewScanner(stdout)
		for sc.Scan() {
			ch <- sc.Text()
		}
		_ = cmd.Wait()
	}()
	return ch, nil
}

func (a *CloudWatchAdapter) FetchLogs(logGroup, since, region string) (string, error) {
	ctx, cancel := defaultContext()
	defer cancel()
	args := []string{"logs", "tail", logGroup, "--since", since}
	if region != "" {
		args = append(args, "--region", region)
	}
	out, err := a.Executor.Run(ctx, "aws", args...)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
