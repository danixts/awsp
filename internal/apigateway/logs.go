package apigateway

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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
	if err := cmd.Start(); err != nil {
		return err
	}
	if follow {
		go func() {
			sc := bufio.NewScanner(os.Stdin)
			for sc.Scan() {
				if strings.TrimSpace(sc.Text()) == "q" {
					_ = cmd.Process.Signal(os.Interrupt)
					return
				}
			}
		}()
	}
	err := cmd.Wait()
	if err != nil {
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
