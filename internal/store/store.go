package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const (
	configDir  = ".config/awsp"
	lastFile   = "last.json"
	favsFile   = "favorites"
)

type LastUsed struct {
	Profile string `json:"profile"`
	Region  string `json:"region"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir), nil
}

func ensureDir() (string, error) {
	dir, err := configPath()
	if err != nil {
		return "", err
	}
	return dir, os.MkdirAll(dir, 0755)
}

func LoadLast() (LastUsed, error) {
	dir, err := configPath()
	if err != nil {
		return LastUsed{}, err
	}
	data, err := os.ReadFile(filepath.Join(dir, lastFile))
	if err != nil {
		return LastUsed{}, err
	}
	var last LastUsed
	if err := json.Unmarshal(data, &last); err != nil {
		return LastUsed{}, err
	}
	return last, nil
}

func SaveLast(profile, region string) error {
	dir, err := ensureDir()
	if err != nil {
		return err
	}
	data, err := json.Marshal(LastUsed{Profile: profile, Region: region})
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, lastFile), data, 0600)
}

func LoadFavorites() ([]string, error) {
	dir, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(dir, favsFile))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out, nil
}

func SaveFavorites(names []string) error {
	dir, err := ensureDir()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, favsFile), []byte(strings.Join(names, "\n")), 0600)
}

func AddFavorite(profile string) error {
	favs, _ := LoadFavorites()
	for _, f := range favs {
		if f == profile {
			return nil
		}
	}
	return SaveFavorites(append(favs, profile))
}

func RemoveFavorite(profile string) error {
	favs, _ := LoadFavorites()
	var out []string
	for _, f := range favs {
		if f != profile {
			out = append(out, f)
		}
	}
	return SaveFavorites(out)
}

func IsFavorite(profile string, favs []string) bool {
	for _, f := range favs {
		if f == profile {
			return true
		}
	}
	return false
}
