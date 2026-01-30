package config

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/danixts/awsp/internal/domain"
	"gopkg.in/ini.v1"
)

const defaultProfile = "default"

type INIProfileLoader struct{}

func (l *INIProfileLoader) LoadProfiles() ([]domain.Profile, error) {
	credPath, configPath, err := awsPaths()
	if err != nil {
		return nil, err
	}

	byName := make(map[string]*domain.Profile)

	credFile, err := ini.Load(credPath)
	if err != nil {
		return nil, err
	}
	for _, sec := range credFile.Sections() {
		name := sec.Name()
		if name == ini.DefaultSection {
			continue
		}
		byName[name] = &domain.Profile{Name: name}
	}

	configFile, err := ini.Load(configPath)
	if err == nil {
		for _, sec := range configFile.Sections() {
			name := sectionToProfileName(sec.Name())
			if name == "" {
				continue
			}
			p, exists := byName[name]
			if !exists {
				p = &domain.Profile{Name: name}
				byName[name] = p
			}
			if key, err := sec.GetKey("region"); err == nil {
				p.Region = key.String()
			}
		}
	}

	names := make([]string, 0, len(byName))
	for name := range byName {
		names = append(names, name)
	}
	sort.Strings(names)

	out := make([]domain.Profile, 0, len(byName))
	for _, name := range names {
		out = append(out, *byName[name])
	}
	return out, nil
}

func (l *INIProfileLoader) CurrentProfile() string {
	if p := os.Getenv("AWS_PROFILE"); p != "" {
		return p
	}
	return defaultProfile
}

func awsPaths() (credPath, configPath string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}
	return filepath.Join(home, ".aws", "credentials"),
		filepath.Join(home, ".aws", "config"), nil
}

func sectionToProfileName(section string) string {
	if section == ini.DefaultSection {
		return ""
	}
	if section == defaultProfile {
		return defaultProfile
	}
	const prefix = "profile "
	if len(section) > len(prefix) && section[:len(prefix)] == prefix {
		return section[len(prefix):]
	}
	return section
}
