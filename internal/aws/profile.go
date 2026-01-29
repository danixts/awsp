package aws

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

const DefaultProfile = "default"

var (
	reCredSection   = regexp.MustCompile(`(?m)^\[([^\]]+)\]`)
	reConfigRegion = regexp.MustCompile(`(?m)^\[(?:profile\s+)?([^\]]+)\][\s\S]*?\n\s*region\s*=\s*(\S+)`)
)

type Profile struct {
	Name   string
	Region string
}

func getPaths() (credPath, configPath string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}
	return filepath.Join(home, ".aws", "credentials"), filepath.Join(home, ".aws", "config"), nil
}

func LoadProfiles() ([]Profile, error) {
	credPath, configPath, err := getPaths()
	if err != nil {
		return nil, err
	}

	byName := make(map[string]*Profile)

	credData, err := os.ReadFile(credPath)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %w", credPath, err)
	}
	for _, m := range reCredSection.FindAllStringSubmatch(string(credData), -1) {
		name := m[1]
		byName[name] = &Profile{Name: name}
	}

	configData, err := os.ReadFile(configPath)
	if err == nil {
		for _, m := range reConfigRegion.FindAllStringSubmatch(string(configData), -1) {
			name, region := m[1], m[2]
			if p, ok := byName[name]; ok {
				p.Region = region
			}
		}
	}

	names := make([]string, 0, len(byName))
	for name := range byName {
		names = append(names, name)
	}
	sort.Strings(names)

	out := make([]Profile, 0, len(byName))
	for _, name := range names {
		out = append(out, *byName[name])
	}
	return out, nil
}

func CurrentProfile() string {
	if p := os.Getenv("AWS_PROFILE"); p != "" {
		return p
	}
	return DefaultProfile
}

func FindProfile(profiles []Profile, name string) (Profile, bool) {
	for _, p := range profiles {
		if p.Name == name {
			return p, true
		}
	}
	return Profile{}, false
}

func OrderWithFavoritesFirst(profiles []Profile, favorites []string) []Profile {
	byName := make(map[string]Profile)
	for _, p := range profiles {
		byName[p.Name] = p
	}
	var out []Profile
	for _, name := range favorites {
		if p, ok := byName[name]; ok {
			out = append(out, p)
			delete(byName, name)
		}
	}
	names := make([]string, 0, len(byName))
	for name := range byName {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		out = append(out, byName[name])
	}
	return out
}
