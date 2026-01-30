package aws

import "github.com/danixts/awsp/internal/adapter/awscli"

var validator = &awscli.ValidatorAdapter{Executor: &awscli.CLIExecutor{}}

func ValidateProfile(profile string) error {
	return validator.Validate(profile)
}
