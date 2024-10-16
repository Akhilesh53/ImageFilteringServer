package models

type GCPSafeSearch struct {
	Adult    string `json:"adult"`
	Medical  string `json:"medical"`
	Racy     string `json:"racy"`
	Spoof    string `json:"spoof"`
	Violence string `json:"violence"`
}

// func to create new instance of GCPSafeSearch
func NewGCPSafeSearch(adult string, medical string, racy string, spoof string, violence string) *GCPSafeSearch {
	return &GCPSafeSearch{
		Adult:    adult,
		Medical:  medical,
		Racy:     racy,
		Spoof:    spoof,
		Violence: violence,
	}
}

// func to get default instance of GCPSafeSearch
func DefaultGCPSafeSearch() *GCPSafeSearch {
	return &GCPSafeSearch{}
}

// get set methods for GCPSafeSearch
func (g *GCPSafeSearch) GetAdult() string {
	return g.Adult
}

func (g *GCPSafeSearch) SetAdult(adult string) *GCPSafeSearch {
	g.Adult = adult
	return g
}

func (g *GCPSafeSearch) GetMedical() string {
	return g.Medical
}

func (g *GCPSafeSearch) SetMedical(medical string) *GCPSafeSearch {
	g.Medical = medical
	return g
}

func (g *GCPSafeSearch) GetRacy() string {
	return g.Racy
}

func (g *GCPSafeSearch) SetRacy(racy string) *GCPSafeSearch {
	g.Racy = racy
	return g
}

func (g *GCPSafeSearch) GetSpoof() string {
	return g.Spoof
}

func (g *GCPSafeSearch) SetSpoof(spoof string) *GCPSafeSearch {
	g.Spoof = spoof
	return g
}

func (g *GCPSafeSearch) GetViolence() string {
	return g.Violence
}

func (g *GCPSafeSearch) SetViolence(violence string) *GCPSafeSearch {
	g.Violence = violence
	return g
}
