package analysis

import (
	"context"
	"sync"
	"time"
)

// getHTMLALinksStatus waits for other goroutines to update required data and
// starts to call each URL to get their status.
func (a *AnalyzerService) getHTMLALinksStatus(ctx context.Context, allLinksAreProcessed <-chan bool) {
	a.logger.InfoF("start getting HTML link statuses | time: %d", time.Now().Unix())

	// Wait for other goroutine to send signal.
	<-allLinksAreProcessed
	a.logger.InfoF("start getting HTML link statuses after signal | time: %d", time.Now().Unix())

	a.result.InaccessibleLinks = 0
	localSync := &sync.WaitGroup{}

	defer func() {
		localSync.Wait()
		a.WaitGroup.Done()
	}()

	// Check internal links.
	for _, link := range a.result.InternalLinks {
		go func() {
			localSync.Add(1)
			err := a.checkLinkStatus(ctx, link)

			if err != nil {
				a.Mutex.Lock()
				a.result.InaccessibleLinks += 1
				a.Mutex.Unlock()
			}

			localSync.Done()
		}()
	}

	// Check external links.
	for _, link := range a.result.ExternalLinks {
		go func() {
			localSync.Add(1)
			err := a.checkLinkStatus(ctx, link)

			if err != nil {
				a.Mutex.Lock()
				a.result.InaccessibleLinks += 1
				a.Mutex.Unlock()
			}

			localSync.Done()
		}()
	}
}

// checkLinkStatus Checks the health of a single link with the help of a third party package
// and updates its struct if the link is unhealthy.
func (a *AnalyzerService) checkLinkStatus(ctx context.Context, link string) error {
	a.logger.InfoF("start getting link status | time: %d", time.Now().Unix())

	err := a.pageParser.GetWebPageStatus(ctx, link)

	if err != nil {
		a.logger.With("link", link).ErrorContext(ctx, err.Error())

		return err
	}

	return nil
}
