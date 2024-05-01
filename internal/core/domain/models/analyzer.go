package models

type Links []string
type HTags map[int]int

type AnalyzerResult struct {
	Title             string
	Version           string
	Hs                HTags
	InternalLinks     Links
	ExternalLinks     Links
	InaccessibleLinks int
	HasLoginForm      bool
}
