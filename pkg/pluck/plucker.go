package pluck

import (
	"bufio"
	"bytes"

	"html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	// external
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	// pp "github.com/sniperkit/colly/plugins/app/debug/pp"

	// internal
	config "github.com/sniperkit/pluck/pkg/config"
	striphtml "github.com/sniperkit/pluck/pkg/striphtml"
)

// Plucker stores the result and the types of things to pluck
type Plucker struct {
	pluckers []pluckUnit
	result   map[string]interface{}
}

type pluckUnit struct {
	config       config.Config
	activators   [][]byte
	patterns     [][]byte
	whitelist    [][]byte
	blacklist    [][]byte
	permanent    int
	maximum      int
	autoSplit    bool
	separator    []byte
	matchMode    []byte
	matchPhrase  []byte
	deactivator  []byte
	finisher     []byte
	captured     [][]byte
	numActivated int
	captureByte  []byte
	captureI     int
	activeI      int
	deactiveI    int
	finisherI    int
	isFinished   bool
}

// New returns a new plucker
// which can later have items added to it
// or can load a config file
// and then can be used to parse.
func New() (*Plucker, error) {
	log.SetLevel(log.WarnLevel)
	p := new(Plucker)
	p.pluckers = []pluckUnit{}
	return p, nil
}

// Verbose toggles debug mode
func (p *Plucker) Verbose(makeVerbose bool) {
	if makeVerbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}

// Configuration returns an array of the current
// Config for each plucker.
func (p *Plucker) Configuration() (c []config.Config) {
	c = make([]config.Config, len(p.pluckers))
	for i, unit := range p.pluckers {
		c[i] = unit.config
	}
	return
}

// Add adds a unit
// to pluck with specified parameters
func (p *Plucker) Add(c config.Config) {
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

	// matchMode
	u.matchMode = []byte(c.Match.Mode)

	// matchPhrase
	u.matchPhrase = []byte(c.Match.Phrase)

	// separator
	u.separator = []byte(c.Match.Separator)

	// autoSplit
	u.autoSplit = c.Match.Split

	// `patterns` is a list of words to extract as a sub-result tree
	u.patterns = make([][]byte, len(c.Patterns))
	for i := range c.Patterns {
		u.patterns[i] = []byte(c.Patterns[i])
	}

	// `whitelist` specifies a word list to check in order to include a match
	u.whitelist = make([][]byte, len(c.Whitelist))
	for i := range c.Whitelist {
		u.whitelist[i] = []byte(c.Whitelist[i])
	}

	// `blacklist` specifies a word list to check in order to exclude a match
	u.blacklist = make([][]byte, len(c.Blacklist))
	for i := range c.Blacklist {
		u.blacklist[i] = []byte(c.Blacklist[i])
	}

	u.captureByte = make([]byte, 100000)
	u.captured = [][]byte{}
	p.pluckers = append(p.pluckers, u)
	log.Infof("Added plucker %+v", c)
}

// Load will load a TOML configuration file of untis
// to pluck with specified parameters
func (p *Plucker) Load(f string) (err error) {
	log.Debugf("load config file at: %s", f)

	var conf *config.Configs
	conf, err = config.NewFromFile(false, false, false, f)
	if err != nil {
		return errors.Wrap(err, "problem opening config file "+f)
	}

	for i := range conf.Pluck {
		var c config.Config
		c.Activators = conf.Pluck[i].Activators
		c.Deactivator = conf.Pluck[i].Deactivator
		c.Finisher = conf.Pluck[i].Finisher
		c.Limit = conf.Pluck[i].Limit
		c.Name = conf.Pluck[i].Name
		c.Permanent = conf.Pluck[i].Permanent
		c.Sanitize = conf.Pluck[i].Sanitize
		c.Maximum = conf.Pluck[i].Maximum
		p.Add(c)
	}

	// Dump config file for dev purpise
	dumpFormats := []string{"yaml", "json", "toml", "xml"}
	dumpNodes := []string{}
	config.Dump(conf, dumpFormats, dumpNodes, "./conf/schema/plucker") // use string slices
	// pp.Println(appConfig)

	return
}

// Load will load a TOML configuration file of untis
// to pluck with specified parameters
func (p *Plucker) LoadTOML(f string) (err error) {
	tomlData, err := ioutil.ReadFile(f)
	if err != nil {
		return errors.Wrap(err, "problem opening "+f)
	}
	log.Debugf("toml string: %s", string(tomlData))
	p.LoadFromString(string(tomlData))
	return
}

// LoadFromString will load a YAML configuration file of untis
// to pluck with specified parameters
func (p *Plucker) LoadFromString(tomlString string) (err error) {
	var conf config.Configs
	_, err = toml.Decode(tomlString, &conf)
	log.Debugf("Loaded toml: %+v", conf)
	for i := range conf.Pluck {
		var c config.Config
		c.Activators = conf.Pluck[i].Activators
		c.Deactivator = conf.Pluck[i].Deactivator
		c.Finisher = conf.Pluck[i].Finisher
		c.Limit = conf.Pluck[i].Limit
		c.Name = conf.Pluck[i].Name
		c.Permanent = conf.Pluck[i].Permanent
		c.Sanitize = conf.Pluck[i].Sanitize
		c.Maximum = conf.Pluck[i].Maximum
		p.Add(c)
	}
	return
}

// PluckString takes a string as input
// and uses the specified parameters and generates
// a map (p.result) with the finished results.
// The streaming can be enabled by setting it to true.
func (p *Plucker) PluckString(s string, stream ...bool) (err error) {
	r := bufio.NewReader(strings.NewReader(s))
	if len(stream) > 0 && stream[0] {
		return p.PluckStream(r)
	}
	return p.Pluck(r)
}

// PluckFile takes a file as input
// and uses the specified parameters and generates
// a map (p.result) with the finished results. The streaming
// can be enabled by setting it to true.
func (p *Plucker) PluckFile(f string, stream ...bool) (err error) {
	r1, err := os.Open(f)
	defer r1.Close()
	if err != nil {
		return
	}
	r := bufio.NewReader(r1)
	if len(stream) > 0 && stream[0] {
		return p.PluckStream(r)
	}
	return p.Pluck(r)
}

// PluckURL takes a URL as input
// and uses the specified parameters and generates
// a map (p.result) with the finished results
func (p *Plucker) PluckURL(url string, stream ...bool) (err error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:52.0) Gecko/20100101 Firefox/52.0")
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	r := bufio.NewReader(resp.Body)
	if len(stream) > 0 && stream[0] {
		return p.PluckStream(r)
	}
	return p.Pluck(r)
}

// Pluck takes a buffered reader stream and
// extracts the text from it. This spawns a thread for
// each plucker and copies the entire buffer to memory,
// so that each plucker works in parallel.
func (p *Plucker) Pluck(r *bufio.Reader) (err error) {
	allBytes, _ := r.ReadBytes(0)
	var wg sync.WaitGroup
	wg.Add(len(p.pluckers))
	for i := 0; i < len(p.pluckers); i++ {
		go func(i int, allBytes []byte) {
			defer wg.Done()
			for _, curByte := range allBytes {
				if p.pluckers[i].numActivated < len(p.pluckers[i].activators) {
					// look for activators
					if curByte == p.pluckers[i].activators[p.pluckers[i].numActivated][p.pluckers[i].activeI] {
						p.pluckers[i].activeI++
						if p.pluckers[i].activeI == len(p.pluckers[i].activators[p.pluckers[i].numActivated]) {
							log.Info(string(curByte), "Activated")
							p.pluckers[i].numActivated++
							p.pluckers[i].activeI = 0
						}
					} else {
						p.pluckers[i].activeI = 0
					}
				} else {
					// add to capture
					p.pluckers[i].captureByte[p.pluckers[i].captureI] = curByte
					p.pluckers[i].captureI++
					// look for deactivators
					if curByte == p.pluckers[i].deactivator[p.pluckers[i].deactiveI] {
						p.pluckers[i].deactiveI++
						if p.pluckers[i].deactiveI == len(p.pluckers[i].deactivator) {
							log.Info(string(curByte), "Deactivated")
							// add capture
							log.Info(string(p.pluckers[i].captureByte[:p.pluckers[i].captureI-len(p.pluckers[i].deactivator)]))
							tempByte := make([]byte, p.pluckers[i].captureI-len(p.pluckers[i].deactivator))
							copy(tempByte, p.pluckers[i].captureByte[:p.pluckers[i].captureI-len(p.pluckers[i].deactivator)])
							if p.pluckers[i].config.Sanitize {
								tempByte = bytes.Replace(tempByte, []byte("\\u003c"), []byte("<"), -1)
								tempByte = bytes.Replace(tempByte, []byte("\\u003e"), []byte(">"), -1)
								tempByte = bytes.Replace(tempByte, []byte("\\u0026"), []byte("&"), -1)
								tempByte = []byte(striphtml.StripTags(html.UnescapeString(string(tempByte))))
							}
							tempByte = bytes.TrimSpace(tempByte)
							if p.pluckers[i].maximum < 1 || len(tempByte) < p.pluckers[i].maximum {
								p.pluckers[i].captured = append(p.pluckers[i].captured, tempByte)
							}
							// reset
							p.pluckers[i].numActivated = p.pluckers[i].permanent
							p.pluckers[i].deactiveI = 0
							p.pluckers[i].captureI = 0
						}
					} else {
						p.pluckers[i].activeI = 0
						p.pluckers[i].deactiveI = 0
					}
				}

				// look for finisher
				if p.pluckers[i].finisher != nil && len(p.pluckers[i].captured) > 0 {
					if curByte == p.pluckers[i].finisher[p.pluckers[i].finisherI] {
						p.pluckers[i].finisherI++
						if p.pluckers[i].finisherI == len(p.pluckers[i].finisher) {
							log.Info(string(curByte), "Finished")
							p.pluckers[i].isFinished = true
						}
					} else {
						p.pluckers[i].finisherI = 0
					}
				}

				if len(p.pluckers[i].captured) == p.pluckers[i].config.Limit {
					p.pluckers[i].isFinished = true
				}
				if p.pluckers[i].isFinished {
					break
				}
			}
			log.Infof("plucker %d finished", i)
		}(i, allBytes)
	}
	wg.Wait()
	p.generateResult()
	return
}

// PluckStream takes a buffered reader stream and streams one
// byte at a time and processes all pluckers serially and
// simultaneously.
func (p *Plucker) PluckStream(r *bufio.Reader) (err error) {
	var finished bool
	for {
		curByte, errRead := r.ReadByte()
		if errRead == io.EOF || finished {
			break
		}
		finished = true
		for i := range p.pluckers {
			if p.pluckers[i].isFinished {
				continue
			}
			finished = false
			if p.pluckers[i].numActivated < len(p.pluckers[i].activators) {
				// look for activators
				if curByte == p.pluckers[i].activators[p.pluckers[i].numActivated][p.pluckers[i].activeI] {
					p.pluckers[i].activeI++
					if p.pluckers[i].activeI == len(p.pluckers[i].activators[p.pluckers[i].numActivated]) {
						log.Info(string(curByte), "Activated")
						p.pluckers[i].numActivated++
						p.pluckers[i].activeI = 0
					}
				} else {
					p.pluckers[i].activeI = 0
				}
			} else {
				// add to capture
				p.pluckers[i].captureByte[p.pluckers[i].captureI] = curByte
				p.pluckers[i].captureI++
				// look for deactivators
				if curByte == p.pluckers[i].deactivator[p.pluckers[i].deactiveI] {
					p.pluckers[i].deactiveI++
					if p.pluckers[i].deactiveI == len(p.pluckers[i].deactivator) {
						log.Info(string(curByte), "Deactivated")
						// add capture
						log.Info(string(p.pluckers[i].captureByte[:p.pluckers[i].captureI-len(p.pluckers[i].deactivator)]))
						tempByte := make([]byte, p.pluckers[i].captureI-len(p.pluckers[i].deactivator))
						copy(tempByte, p.pluckers[i].captureByte[:p.pluckers[i].captureI-len(p.pluckers[i].deactivator)])
						if p.pluckers[i].config.Sanitize {
							tempByte = bytes.Replace(tempByte, []byte("\\u003c"), []byte("<"), -1)
							tempByte = bytes.Replace(tempByte, []byte("\\u003e"), []byte(">"), -1)
							tempByte = bytes.Replace(tempByte, []byte("\\u0026"), []byte("&"), -1)
							tempByte = []byte(striphtml.StripTags(html.UnescapeString(string(tempByte))))
						}
						tempByte = bytes.TrimSpace(tempByte)
						p.pluckers[i].captured = append(p.pluckers[i].captured, tempByte)
						// reset
						p.pluckers[i].numActivated = p.pluckers[i].permanent
						p.pluckers[i].deactiveI = 0
						p.pluckers[i].captureI = 0
					}
				} else {
					p.pluckers[i].activeI = 0
					p.pluckers[i].deactiveI = 0
				}
			}

			// look for finisher
			if p.pluckers[i].finisher != nil {
				if curByte == p.pluckers[i].finisher[p.pluckers[i].finisherI] {
					p.pluckers[i].finisherI++
					if p.pluckers[i].finisherI == len(p.pluckers[i].finisher) {
						log.Info(string(curByte), "Finished")
						p.pluckers[i].isFinished = true
					}
				} else {
					p.pluckers[i].finisherI = 0
				}
			}

			if len(p.pluckers[i].captured) == p.pluckers[i].config.Limit {
				p.pluckers[i].isFinished = true
			}
		}
	}
	p.generateResult()
	return
}

func (p *Plucker) generateResult() {
	p.result = make(map[string]interface{})
	for i := range p.pluckers {
		if len(p.pluckers[i].captured) == 0 {
			p.result[p.pluckers[i].config.Name] = ""
		} else if len(p.pluckers[i].captured) == 1 {
			p.result[p.pluckers[i].config.Name] = string(p.pluckers[i].captured[0])
		} else {
			results := make([]string, len(p.pluckers[i].captured))
			for j, r := range p.pluckers[i].captured {
				results[j] = string(r)
			}
			if len(results) == 0 {
				p.result[p.pluckers[i].config.Name] = ""
			} else {
				p.result[p.pluckers[i].config.Name] = results
			}
		}
	}
}

// Result returns the raw result
func (p *Plucker) Result() map[string]interface{} {
	return p.result
}
