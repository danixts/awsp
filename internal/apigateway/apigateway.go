package apigateway

import (
	"github.com/danixts/awsp/internal/adapter/awscli"
	"github.com/danixts/awsp/internal/domain"
)

// Type aliases for backward compatibility with existing TUI code.
type RestAPI = domain.RestAPI
type ResourceWithMethods = domain.ResourceWithMethods
type Endpoint = domain.Endpoint
type TimeRange = domain.TimeRange

var DefaultTimeRanges = domain.DefaultTimeRanges

var gw = &awscli.GatewayAdapter{Executor: &awscli.CLIExecutor{}}

func GetRestAPIs() ([]RestAPI, error) {
	return gw.ListAPIs()
}

func GetResources(restAPIID string) ([]ResourceWithMethods, error) {
	return gw.ListResources(restAPIID)
}

func EndpointsFromResources(resources []ResourceWithMethods) []Endpoint {
	return domain.EndpointsFromResources(resources)
}

func GetIntegrationLambdaLogGroup(restAPIID, resourceID, httpMethod string) (string, error) {
	return gw.GetIntegrationLogGroup(restAPIID, resourceID, httpMethod)
}
