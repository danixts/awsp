package domain

type Profile struct {
	Name   string
	Region string
}

type Region struct {
	Code string
	Name string
}

type LastUsed struct {
	Profile string `json:"profile"`
	Region  string `json:"region"`
}
