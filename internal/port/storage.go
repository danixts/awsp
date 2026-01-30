package port

import "github.com/danixts/awsp/internal/domain"

type ConfigStore interface {
	LoadLast() (domain.LastUsed, error)
	SaveLast(profile, region string) error
	LoadFavorites() ([]string, error)
	SaveFavorites(names []string) error
	AddFavorite(profile string) error
	RemoveFavorite(profile string) error
	LoadTheme() (string, error)
	SaveTheme(name string) error
}
