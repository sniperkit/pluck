package dashboard

/*
	Notes:
	- Create an interface for swapping easily from termui, tview/tcell gocui libraries
*/

// A TermUI interface is used to swap between consule ui libraries
type TermUI interface {
	Dashboard(stopTheUI, stopTheCrawler chan bool)
}

func New() {}
