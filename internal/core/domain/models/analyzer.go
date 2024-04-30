package models

type HTags map[int]int

type AnalyzerResult struct {
	Title   string
	Version string
	Hs      HTags
}
