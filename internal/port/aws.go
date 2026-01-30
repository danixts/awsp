package port

import "github.com/danixts/awsp/internal/domain"

type ProfileLoader interface {
	LoadProfiles() ([]domain.Profile, error)
	CurrentProfile() string
}

type CredentialValidator interface {
	Validate(profileName string) error
}
