package aws

import "github.com/danixts/awsp/internal/domain"

// Region is an alias for domain.Region for backward compatibility.
type Region = domain.Region

var CommonRegions = []Region{
	// North America
	{Code: "us-east-1", Name: "US East (N. Virginia)"},
	{Code: "us-east-2", Name: "US East (Ohio)"},
	{Code: "us-west-1", Name: "US West (N. California)"},
	{Code: "us-west-2", Name: "US West (Oregon)"},
	{Code: "ca-central-1", Name: "Canada (Central)"},
	{Code: "ca-west-1", Name: "Canada West (Calgary)"},
	// Europe
	{Code: "eu-west-1", Name: "Europe (Ireland)"},
	{Code: "eu-west-2", Name: "Europe (London)"},
	{Code: "eu-west-3", Name: "Europe (Paris)"},
	{Code: "eu-central-1", Name: "Europe (Frankfurt)"},
	{Code: "eu-central-2", Name: "Europe (Zurich)"},
	{Code: "eu-north-1", Name: "Europe (Stockholm)"},
	{Code: "eu-south-1", Name: "Europe (Milan)"},
	{Code: "eu-south-2", Name: "Europe (Spain)"},
	// Asia Pacific
	{Code: "ap-southeast-1", Name: "Asia Pacific (Singapore)"},
	{Code: "ap-southeast-2", Name: "Asia Pacific (Sydney)"},
	{Code: "ap-southeast-3", Name: "Asia Pacific (Jakarta)"},
	{Code: "ap-southeast-4", Name: "Asia Pacific (Melbourne)"},
	{Code: "ap-northeast-1", Name: "Asia Pacific (Tokyo)"},
	{Code: "ap-northeast-2", Name: "Asia Pacific (Seoul)"},
	{Code: "ap-northeast-3", Name: "Asia Pacific (Osaka)"},
	{Code: "ap-south-1", Name: "Asia Pacific (Mumbai)"},
	{Code: "ap-south-2", Name: "Asia Pacific (Hyderabad)"},
	{Code: "ap-east-1", Name: "Asia Pacific (Hong Kong)"},
	// South America
	{Code: "sa-east-1", Name: "South America (Sao Paulo)"},
	// Middle East
	{Code: "me-south-1", Name: "Middle East (Bahrain)"},
	{Code: "me-central-1", Name: "Middle East (UAE)"},
	// Africa
	{Code: "af-south-1", Name: "Africa (Cape Town)"},
	// Israel
	{Code: "il-central-1", Name: "Israel (Tel Aviv)"},
}

func RegionListForSelector(profileRegion string) []Region {
	if profileRegion == "" {
		return CommonRegions
	}
	for _, r := range CommonRegions {
		if r.Code == profileRegion {
			return CommonRegions
		}
	}
	def := Region{Code: profileRegion, Name: "(profile default)"}
	return append([]Region{def}, CommonRegions...)
}
