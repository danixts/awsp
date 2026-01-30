package domain

import "errors"

var (
	ErrProfileNotFound      = errors.New("profile not found")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrAWSCLINotFound       = errors.New("aws CLI not found in PATH")
	ErrNoProfiles           = errors.New("no profiles configured")
	ErrNotLambdaIntegration = errors.New("integration is not Lambda-backed or URI format unknown")
	ErrNoAPIsFound          = errors.New("no API Gateways found")
	ErrNoEndpoints          = errors.New("no endpoints with methods")
)
