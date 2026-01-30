package apigateway

import (
	"context"

	"github.com/danixts/awsp/internal/adapter/awscli"
)

var cw = &awscli.CloudWatchAdapter{Executor: &awscli.CLIExecutor{}}

func TailLogsStream(ctx context.Context, logGroup, since, region string) <-chan string {
	ch, err := cw.TailLogs(ctx, logGroup, since, region)
	if err != nil {
		errCh := make(chan string, 1)
		errCh <- "Error: " + err.Error()
		close(errCh)
		return errCh
	}
	return ch
}

func TailLogsContent(logGroup, since, region string) (string, error) {
	return cw.FetchLogs(logGroup, since, region)
}
