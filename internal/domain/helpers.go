package domain

import "sort"

func EndpointsFromResources(resources []ResourceWithMethods) []Endpoint {
	var out []Endpoint
	for _, r := range resources {
		for _, m := range r.Methods {
			out = append(out, Endpoint{ResourceID: r.ResourceID, Path: r.Path, Method: m})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Path != out[j].Path {
			return out[i].Path < out[j].Path
		}
		return out[i].Method < out[j].Method
	})
	return out
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

// FindProfile returns the profile with the given name.
func FindProfile(profiles []Profile, name string) (Profile, bool) {
	for _, p := range profiles {
		if p.Name == name {
			return p, true
		}
	}
	return Profile{}, false
}

func IsFavorite(profile string, favs []string) bool {
	for _, f := range favs {
		if f == profile {
			return true
		}
	}
	return false
}
