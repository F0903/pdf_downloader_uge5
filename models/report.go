package models

type Report struct {
	Id                   string
	Name                 string
	PrimaryDownloadLink  string
	FallbackDownloadLink string
}

// Implement Downloadable

// Returns urls in order of importance
func (report *Report) GetDownloadableURLs() []string {
	return []string{
		report.PrimaryDownloadLink,
		report.FallbackDownloadLink,
	}
}
