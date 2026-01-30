package aws

import (
	"github.com/danixts/awsp/internal/adapter/config"
	"github.com/danixts/awsp/internal/domain"
)

const DefaultProfile = "default"

// Profile is an alias for domain.Profile for backward compatibility.
type Profile = domain.Profile

var loader = &config.INIProfileLoader{}

func LoadProfiles() ([]Profile, error) {
	return loader.LoadProfiles()
}

func CurrentProfile() string {
	return loader.CurrentProfile()
}

func FindProfile(profiles []Profile, name string) (Profile, bool) {
	return domain.FindProfile(profiles, name)
}

func OrderWithFavoritesFirst(profiles []Profile, favorites []string) []Profile {
	return domain.OrderWithFavoritesFirst(profiles, favorites)
}
