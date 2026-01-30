package store

import (
	"github.com/danixts/awsp/internal/adapter/store"
	"github.com/danixts/awsp/internal/domain"
)

// LastUsed is an alias for domain.LastUsed for backward compatibility.
type LastUsed = domain.LastUsed

var fs = &store.FileStore{}

func LoadLast() (LastUsed, error) {
	return fs.LoadLast()
}

func SaveLast(profile, region string) error {
	return fs.SaveLast(profile, region)
}

func LoadFavorites() ([]string, error) {
	return fs.LoadFavorites()
}

func SaveFavorites(names []string) error {
	return fs.SaveFavorites(names)
}

func AddFavorite(profile string) error {
	return fs.AddFavorite(profile)
}

func RemoveFavorite(profile string) error {
	return fs.RemoveFavorite(profile)
}

func IsFavorite(profile string, favs []string) bool {
	return domain.IsFavorite(profile, favs)
}

func LoadGatewayTheme() (string, error) {
	return fs.LoadTheme()
}

func SaveGatewayTheme(themeName string) error {
	return fs.SaveTheme(themeName)
}
