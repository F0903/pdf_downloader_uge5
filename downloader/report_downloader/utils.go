package report_downloader

func CountSuccesfulReportDownloads(results []*ReportDownloadResult) int {
	counter := 0
	for _, result := range results {
		if !result.State.IsDone() {
			continue
		}
		counter += 1
	}
	return counter
}
