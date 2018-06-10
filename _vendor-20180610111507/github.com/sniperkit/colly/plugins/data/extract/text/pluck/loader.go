package plucker

import (
	"io/ioutil"
	"strconv"

	"github.com/pkg/errors"

	cfg "github.com/sniperkit/colly/plugins/data/extract/text/pluck/config"
)

type Configs struct {
	Pluck []cfg.Config
}

// NewWithConfig instanciate a new plucker with a plucker config structure
func NewWithConfig(configs ...cfg.Config) (*Plucker, error) {
	p, err := New()
	if err != nil {
		return nil, err
	}
	for _, config := range configs {
		p.Add(config)
	}
	return p, nil
}

// Load will load a YAML configuration file of untis
// to pluck with specified parameters
func (p *Plucker) Load(f string) (err error) {
	config, err := ioutil.ReadFile(f)
	if err != nil {
		return errors.Wrap(err, errFailedToOpen+f)
	}
	// log.Debugf("toml string: %s", string(tomlData))
	p.LoadFromString(string(config))
	return
}

// Add adds a unit
// to pluck with specified parameters
// Nb. Validate rules ?!!!!
func (p *Plucker) Add(c cfg.Config) {
	var u pluckUnit
	u.config = c
	if u.config.Limit == 0 {
		u.config.Limit = -1
	}
	if u.config.Name == "" {
		u.config.Name = strconv.Itoa(len(p.pluckers))
	}
	u.activators = make([][]byte, len(c.Activators))
	for i := range c.Activators {
		u.activators[i] = []byte(c.Activators[i])
	}
	u.permanent = c.Permanent
	u.deactivator = []byte(c.Deactivator)
	if len(c.Finisher) > 0 {
		u.finisher = []byte(c.Finisher)
	} else {
		u.finisher = nil
	}
	u.maximum = -1
	if c.Maximum > 0 {
		u.maximum = c.Maximum
	}
	u.captureByte = make([]byte, 100000)
	u.captured = [][]byte{}
	p.pluckers = append(p.pluckers, u)
	// log.Infof("Added plucker %+v", c)
}

// Config returns an array of the current setup of each plucking unit.
func (p *Plucker) Config() (c []cfg.Config) {
	c = make([]cfg.Config, len(p.pluckers))
	for i, unit := range p.pluckers {
		c[i] = unit.config
	}
	return
}
