package tablib

import (
	"sync"
)

// Sheet represents a sheet in a Databook, holding a title (if any) and a dataset.
type Sheet struct {
	title   string
	dataset *Dataset
	lock    *sync.RWMutex
	wg      *sync.WaitGroup
}

// Title return the title of the sheet.
func (s Sheet) Title() string {
	// s.lock.RLock()
	// defer s.lock.RUnlock()
	return s.title
}

// Dataset returns the dataset of the sheet.
func (s Sheet) Dataset() *Dataset {
	// s.lock.RLock()
	// defer s.lock.RUnlock()

	return s.dataset
}

// Databook represents a Databook which is an array of sheets.
type Databook struct {
	sheets map[string]Sheet
	lock   *sync.RWMutex
	wg     *sync.WaitGroup
}

// NewDatabook constructs a new Databook.
func NewDatabook() *Databook {
	return &Databook{
		sheets: make(map[string]Sheet),
		lock:   &sync.RWMutex{},
		wg:     &sync.WaitGroup{},
	}
}

// Sheets returns the sheets in the Databook.
func (d *Databook) Sheets() map[string]Sheet {
	return d.sheets
}

// Sheet returns the sheet with a specific title.
func (d *Databook) Sheet(title string) Sheet {
	return d.sheets[title]
}

// AddSheet adds a sheet to the Databook.
func (d *Databook) AddSheet(title string, dataset *Dataset) {
	d.sheets[title] = Sheet{
		title:   title,
		dataset: dataset,
		lock:    &sync.RWMutex{},
		wg:      &sync.WaitGroup{},
	}
}

// Size returns the number of sheets in the Databook.
func (d *Databook) Size() int {
	return len(d.sheets)
}

// Wipe removes all Dataset objects from the Databook.
func (d *Databook) Wipe() {
	for k := range d.sheets {
		delete(d.sheets, k)
	}
}
