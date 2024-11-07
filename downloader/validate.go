package downloader

import (
	"fmt"
	"sync"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func ValidateDownloadResult(result *DownloadResult) {
	state := result.State
	if !state.IsDone() {
		return
	}

	filePath := state.writtenPath
	err := api.ValidateFile(filePath, model.NewDefaultConfiguration())
	if err == nil {
		return
	}

	state.SetError(fmt.Errorf("could not validate PDF: %w", err))
}

// Currently not used. It makes more sense to validate each download individually the moment they are downloaded.
func ValidateDownloadResults(results []*DownloadResult) {
	var wg sync.WaitGroup
	p := mpb.New(
		mpb.WithWaitGroup(&wg),
		mpb.WithAutoRefresh(),
	)

	progressBar := p.AddBar(
		int64(len(results)),
		mpb.PrependDecorators(
			decor.Name("Validating documents...", decor.WC{C: decor.DindentRight | decor.DextraSpace}),
		),
		mpb.AppendDecorators(
			decor.AverageETA(decor.ET_STYLE_GO, decor.WC{C: decor.DindentRight | decor.DextraSpace}),
			decor.Percentage(),
		),
		mpb.BarRemoveOnComplete(),
	)

	for _, result := range results {
		wg.Add(1)
		go func() {
			defer func() {
				progressBar.Increment()
				wg.Done()
			}()
			ValidateDownloadResult(result)
		}()
	}

	p.Wait()
}
