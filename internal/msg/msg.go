package msg

const (
	ErrLoadProfiles    = "Error: %v"
	ErrNoProfiles      = "No profiles found in ~/.aws/credentials"
	ErrProfileNotFound = "Profile '%s' not found"
	Canceled           = "Canceled"
	ExportProfile       = "export AWS_PROFILE=%s\n"
	ExportRegion        = "export AWS_REGION=%s\nexport AWS_DEFAULT_REGION=%s\n"
	ExportSDKLoadConfig  = "export AWS_SDK_LOAD_CONFIG=1\n"
	ExportClearCreds = "unset AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_SESSION_TOKEN 2>/dev/null\n"
	HintApply          = "# Apply in this shell: eval $(aws-profile)\n"
	CurrentProfile     = "Current profile: %s\n"
	ListItemWithReg    = "  %s (%s)\n"
	ListItem           = "  %s\n"
	ValidateFailed     = "Profile validation failed: %v\n"
	FavoriteAdded      = "Added '%s' to favorites\n"
	FavoriteRemoved   = "Removed '%s' from favorites\n"
	FavoritesList      = "Favorites:\n"
	NoFavorites       = "No favorites set\n"
)
