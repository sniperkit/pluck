package advanced

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/sniperkit/colly/pkg"
)

type ProxyPool interface {
	SetPool() error
}

type Proxy struct {
	p     string `required:'true'`
	ready bool   `default:'false'`
	isTor bool   `default:'false'`
	lock  *sync.RWMutex
	err   error
}

// Constructor
func New(ps string) *Proxy {
	return &Proxy{
		p:    ps,
		lock: &sync.RWMutex{},
	}
}

// Constructor
func NewWithConfig(c *Config) *Proxy {
	ps := fmt.Sprintf("%s://%s:%d", c.Protocol, c.Host, c.Port)
	p := &Proxy{
		p:    ps,
		lock: &sync.RWMutex{},
	}
	if c.Port == 9050 {
		p.isTor = true
	}
	// To do:
	// - check if healthy
	// - check if valid
	// - fetch proxy list...
	return p
}

func (p *Proxy) IsReady() bool {
	return bool(p.ready)
}

func (p *Proxy) IsOnion() bool {
	return bool(p.isTor)
}

func (p *Proxy) IsHealthy() (ok bool) {
	if p.err == nil {
		ok = true
	}
	return
}

func (p *Proxy) String() string {
	return string(p.p)
}

type roundRobinSwitcher struct {
	proxyURLs []*url.URL
	index     uint32
}

func (r *roundRobinSwitcher) GetProxy(_ *http.Request) (*url.URL, error) {
	u := r.proxyURLs[r.index%uint32(len(r.proxyURLs))]
	atomic.AddUint32(&r.index, 1)
	return u, nil
}

// RoundRobinProxySwitcher creates a proxy switcher function which rotates
// ProxyURLs on every request.
// The proxy type is determined by the URL scheme. "http", "https"
// and "socks5" are supported. If the scheme is empty,
// "http" is assumed.
func RoundRobinProxySwitcher(ProxyURLs ...string) (colly.ProxyFunc, error) {
	urls := make([]*url.URL, len(ProxyURLs))
	for i, u := range ProxyURLs {
		parsedU, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		urls[i] = parsedU
	}
	return (&roundRobinSwitcher{urls, 0}).GetProxy, nil
}
