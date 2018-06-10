package cui

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	cfg "github.com/sniperkit/colly/pkg/config"
)

type TermUI struct {
	Config *cfg.Config
	app    *tview.Application
	list   *tview.List
	txt    *tview.TextView
	event  *tcell.EventKey
}

func Dashboard(stopTheUI, stopTheCrawler chan bool) {

	// stopTheCrawler <- true

	// stop when the crawler is done
	go func() {
		select {
		// case <-pauseTheUI:
		case <-stopTheUI:
			// wait 10 seconds before closing the ui
			time.Sleep(time.Second * 10)
			// termui.StopLoop()
		}
	}()

}
