package port

import "github.com/danixts/awsp/internal/domain"

type GatewayService interface {
	ListAPIs() ([]domain.RestAPI, error)
	ListResources(apiID string) ([]domain.ResourceWithMethods, error)
	GetIntegrationLogGroup(apiID, resourceID, method string) (string, error)
}
