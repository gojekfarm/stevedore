package stevedore

// ChartRepository represents a chart repository
type ChartRepository interface {
	DownloadIndexFile() (string, error)
}
