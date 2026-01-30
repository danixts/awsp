package domain

import "errors"

var ErrNotLambdaIntegration = errors.New("integration is not Lambda-backed or URI format unknown")
