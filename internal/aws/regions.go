package aws

type Region struct {
	Code string
	Name string
}

var CommonRegions = []Region{
	{Code: "us-east-1", Name: "US East (N. Virginia)"},
	{Code: "us-east-2", Name: "US East (Ohio)"},
	{Code: "us-west-1", Name: "US West (N. California)"},
	{Code: "us-west-2", Name: "US West (Oregon)"},
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
