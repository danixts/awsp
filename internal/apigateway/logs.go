package apigateway

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func TailLogsContent(logGroup, since, region string) (string, error) {
	args := []string{"logs", "tail", logGroup, "--since", since}
	if region != "" {
		args = append(args, "--region", region)
	}
	cmd := exec.Command("aws", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("logs tail: %w", err)
	}
	return out.String(), nil
}

func TailLogsStream(logGroup, since, region string) <-chan string {
	ch := make(chan string, 100)
	go func() {
		defer close(ch)
		args := []string{"logs", "tail", logGroup, "--since", since, "--follow"}
		if region != "" {
			args = append(args, "--region", region)
		}
		cmd := exec.Command("aws", args...)
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
	return ch
}

func TailLogs(logGroup, since, region string, follow bool) error {
	args := []string{"logs", "tail", logGroup, "--since", since}
	if region != "" {
		args = append(args, "--region", region)
	}
	if follow {
		args = append(args, "--follow")
	}
	cmd := exec.Command("aws", args...)
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("logs tail exit %d", ee.ExitCode())
		}
		return err
	}
	return nil
}

type TimeRange struct {
	Label string
	Since string
}

var DefaultTimeRanges = []TimeRange{
	{Label: "Last 5 minutes", Since: "5m"},
	{Label: "Last 15 minutes", Since: "15m"},
	{Label: "Last 30 minutes", Since: "30m"},
	{Label: "Last 1 hour", Since: "1h"},
	{Label: "Last 3 hours", Since: "3h"},
	{Label: "Last 12 hours", Since: "12h"},
	{Label: "Last 24 hours (today)", Since: "24h"},
	{Label: "Last 7 days", Since: "7d"},
}
