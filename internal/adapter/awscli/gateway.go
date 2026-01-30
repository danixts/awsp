package awscli

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"

	"github.com/danixts/awsp/internal/domain"
)

var lambdaURIRegex = regexp.MustCompile(`function:([^/]+)/invocations`)

type GatewayAdapter struct {
	Executor CommandExecutor
}

type restAPIsResponse struct {
	Items []domain.RestAPI `json:"items"`
}

type resourceItem struct {
	ID              string                 `json:"id"`
	Path            string                 `json:"path"`
	ResourceMethods map[string]interface{} `json:"resourceMethods"`
}

type resourcesResponse struct {
	Items []resourceItem `json:"items"`
}

func (a *GatewayAdapter) ListAPIs() ([]domain.RestAPI, error) {
	ctx, cancel := defaultContext()
	defer cancel()
	out, err := a.Executor.Run(ctx, "aws", "apigateway", "get-rest-apis", "--output", "json")
	if err != nil {
		return nil, fmt.Errorf("list APIs: %w", err)
	}
	var resp restAPIsResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, fmt.Errorf("parse get-rest-apis: %w", err)
	}
	return resp.Items, nil
}

func (a *GatewayAdapter) ListResources(apiID string) ([]domain.ResourceWithMethods, error) {
	ctx, cancel := defaultContext()
	defer cancel()
	out, err := a.Executor.Run(ctx, "aws", "apigateway", "get-resources",
		"--rest-api-id", apiID, "--output", "json")
	if err != nil {
		return nil, fmt.Errorf("list resources: %w", err)
	}
	var resp resourcesResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, fmt.Errorf("parse get-resources: %w", err)
	}
	var result []domain.ResourceWithMethods
	for _, r := range resp.Items {
		if len(r.ResourceMethods) == 0 {
			continue
		}
		methods := make([]string, 0, len(r.ResourceMethods))
		for m := range r.ResourceMethods {
			methods = append(methods, m)
		}
		sort.Strings(methods)
		result = append(result, domain.ResourceWithMethods{
			ResourceID: r.ID,
			Path:       r.Path,
			Methods:    methods,
		})
	}
	return result, nil
}

func (a *GatewayAdapter) GetIntegrationLogGroup(apiID, resourceID, method string) (string, error) {
	ctx, cancel := defaultContext()
	defer cancel()
	out, err := a.Executor.Run(ctx, "aws", "apigateway", "get-integration",
		"--rest-api-id", apiID,
		"--resource-id", resourceID,
		"--http-method", method,
		"--output", "json")
	if err != nil {
		return "", fmt.Errorf("get-integration: %w", err)
	}
	var resp struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(out, &resp); err != nil {
		return "", fmt.Errorf("parse get-integration: %w", err)
	}
	matches := lambdaURIRegex.FindStringSubmatch(resp.URI)
	if len(matches) < 2 {
		return "", domain.ErrNotLambdaIntegration
	}
	return "/aws/lambda/" + matches[1], nil
}
