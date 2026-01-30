package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/danixts/awsp/internal/domain"
)

const (
	configDir    = ".config/awsp"
	lastFile     = "last.json"
	favsFile     = "favorites"
	gatewayTheme = "gateway_theme.json"
)

type FileStore struct{}

func (s *FileStore) configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir), nil
}

func (s *FileStore) ensureDir() (string, error) {
	dir, err := s.configPath()
	if err != nil {
		return "", err
	}
	return dir, os.MkdirAll(dir, 0755)
}

func (s *FileStore) LoadLast() (domain.LastUsed, error) {
	dir, err := s.configPath()
	if err != nil {
		return domain.LastUsed{}, err
	}
	data, err := os.ReadFile(filepath.Join(dir, lastFile))
	if err != nil {
		return domain.LastUsed{}, err
	}
	var last domain.LastUsed
	if err := json.Unmarshal(data, &last); err != nil {
		return domain.LastUsed{}, err
	}
	return last, nil
}

func (s *FileStore) SaveLast(profile, region string) error {
	dir, err := s.ensureDir()
	if err != nil {
		return err
	}
	data, err := json.Marshal(domain.LastUsed{Profile: profile, Region: region})
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, lastFile), data, 0600)
}

func (s *FileStore) LoadFavorites() ([]string, error) {
	dir, err := s.configPath()
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

func (s *FileStore) SaveFavorites(names []string) error {
	dir, err := s.ensureDir()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, favsFile), []byte(strings.Join(names, "\n")), 0600)
}

func (s *FileStore) AddFavorite(profile string) error {
	favs, _ := s.LoadFavorites()
	for _, f := range favs {
		if f == profile {
			return nil
		}
	}
	return s.SaveFavorites(append(favs, profile))
}

func (s *FileStore) RemoveFavorite(profile string) error {
	favs, _ := s.LoadFavorites()
	var out []string
	for _, f := range favs {
		if f != profile {
			out = append(out, f)
		}
	}
	return s.SaveFavorites(out)
}

type gatewayThemeData struct {
	Theme string `json:"theme"`
}

func (s *FileStore) LoadTheme() (string, error) {
	dir, err := s.configPath()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(filepath.Join(dir, gatewayTheme))
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	var g gatewayThemeData
	if err := json.Unmarshal(data, &g); err != nil {
		return "", err
	}
	return g.Theme, nil
}

func (s *FileStore) SaveTheme(themeName string) error {
	dir, err := s.ensureDir()
	if err != nil {
		return err
	}
	data, err := json.Marshal(gatewayThemeData{Theme: themeName})
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, gatewayTheme), data, 0600)
}
