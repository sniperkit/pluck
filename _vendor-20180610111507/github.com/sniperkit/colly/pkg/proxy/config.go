package proxy

const (
	REMOTE_PROXY_LIST_URL string = "https://proxy-list.org/english/index.php"
)

var (
	defaultTorHost     string = "127.0.0.1"
	defaultTorPort     string = "9050"
	defaultTorProtocol string = "socks5"
)

type Config struct {
	Protocol string `socks5`
	Host     string `default:'127.0.0.1'`
	Port     int    `default:'9050'`
	Fetch    bool   `default:'true'`
}
