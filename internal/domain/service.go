package domain

type RestAPI struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ResourceWithMethods struct {
	ResourceID string
	Path       string
	Methods    []string
}

type Endpoint struct {
	ResourceID string
	Path       string
	Method     string
}

type LogGroup struct {
	Name            string `json:"logGroupName"`
	StoredBytes     int64  `json:"storedBytes"`
	RetentionInDays int    `json:"retentionInDays"`
}

type TimeRange struct {
	Label string
	Since string
}

var DefaultTimeRanges = []TimeRange{
	{Label: "Last 5 minutes", Since: "5m"},
	{Label: "Last 15 minutes", Since: "15m"},
	{Label: "Last 30 minutes", Since: "30m"},
	{Label: "Last 1 hour", Since: "1h"},
	{Label: "Last 3 hours", Since: "3h"},
	{Label: "Last 12 hours", Since: "12h"},
	{Label: "Last 24 hours (today)", Since: "24h"},
	{Label: "Last 7 days", Since: "7d"},
}
