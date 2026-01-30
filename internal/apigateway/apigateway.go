package apigateway

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
)

type RestAPI struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type restAPIsResponse struct {
	Items []RestAPI `json:"items"`
}

type Resource struct {
	ID              string                 `json:"id"`
	Path            string                 `json:"path"`
	ResourceMethods map[string]interface{} `json:"resourceMethods"`
}

type resourcesResponse struct {
	Items []Resource `json:"items"`
}

func GetRestAPIs() ([]RestAPI, error) {
	out, err := exec.Command("aws", "apigateway", "get-rest-apis", "--output", "json").Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("aws apigateway: %s", string(ee.Stderr))
		}
		return nil, err
	}
	var resp restAPIsResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, fmt.Errorf("parse get-rest-apis: %w", err)
	}
	return resp.Items, nil
}

func GetResources(restAPIID string) ([]ResourceWithMethods, error) {
	out, err := exec.Command("aws", "apigateway", "get-resources",
		"--rest-api-id", restAPIID,
		"--output", "json").Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("aws apigateway: %s", string(ee.Stderr))
		}
		return nil, err
	}
	var resp resourcesResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, fmt.Errorf("parse get-resources: %w", err)
	}
	var result []ResourceWithMethods
	for _, r := range resp.Items {
		if len(r.ResourceMethods) == 0 {
			continue
		}
		methods := make([]string, 0, len(r.ResourceMethods))
		for m := range r.ResourceMethods {
			methods = append(methods, m)
		}
		sort.Strings(methods)
		result = append(result, ResourceWithMethods{ResourceID: r.ID, Path: r.Path, Methods: methods})
	}
	return result, nil
}

type ResourceWithMethods struct {
	ResourceID string
	Path       string
	Methods    []string
}

type Endpoint struct {
	ResourceID string
	Path       string
	Method     string
}

func EndpointsFromResources(resources []ResourceWithMethods) []Endpoint {
	var out []Endpoint
	for _, r := range resources {
		for _, m := range r.Methods {
			out = append(out, Endpoint{ResourceID: r.ResourceID, Path: r.Path, Method: m})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Path != out[j].Path {
			return out[i].Path < out[j].Path
		}
		return out[i].Method < out[j].Method
	})
	return out
}

var lambdaURIRegex = regexp.MustCompile(`function:([^/]+)/invocations`)

func GetIntegrationLambdaLogGroup(restAPIID, resourceID, httpMethod string) (logGroup string, err error) {
	out, err := exec.Command("aws", "apigateway", "get-integration",
		"--rest-api-id", restAPIID,
		"--resource-id", resourceID,
		"--http-method", httpMethod,
		"--output", "json").Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("get-integration: %s", string(ee.Stderr))
		}
		return "", err
	}
	var resp struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(out, &resp); err != nil {
		return "", fmt.Errorf("parse get-integration: %w", err)
	}
	matches := lambdaURIRegex.FindStringSubmatch(resp.URI)
	if len(matches) < 2 {
		return "", fmt.Errorf("integration is not Lambda or URI format unknown")
	}
	return "/aws/lambda/" + matches[1], nil
}
