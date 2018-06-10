package main

import (
	"sync"

	metric "github.com/sniperkit/colly/pkg/metric"
	tui "github.com/sniperkit/colly/plugins/app/dashboard/tui/termui"
)

var (
	stopTheUI                    chan bool
	stopTheCrawler               chan bool
	allURLsHaveBeenVisited       chan bool
	allStatisticsHaveBeenUpdated chan bool
	uiWaitGroup                  = &sync.WaitGroup{}
)

func initStatsCollector() {
	collectorStats = metric.NewStatsCollector()
}

func initDashboard() {
	collectorStats = metric.NewStatsCollector()

	stopTheUI = make(chan bool)
	collectorResponseMetrics = make(chan metric.Response)

	// uiWaitGroup.Add(1)
	go func() {
		tui.Dashboard(collectorStats, stopTheUI, stopTheCrawler)
		// uiWaitGroup.Done()
	}()
	// uiWaitGroup.Wait()

}

func updateDashboard() {
	go func() {
		for {
			select {
			case <-stopTheCrawler:
				stopTheUI <- true
				return

			case <-stopTheUI:
				stopTheCrawler <- true

			case snapshot := <-collectorResponseMetrics:
				if collectorStats != nil {
					collectorStats.UpdateStatistics(snapshot)
				}
			}
		}
	}()
}
