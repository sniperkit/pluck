package proxy

// Config
type Config struct {

	// Enabled
	Enabled bool `default:"false" flag:"with-backoff" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled" csv:"Enabled"`

	// FetchRemoteList
	FetchRemoteList bool `default:"true" flag:"-fetch-remote-list" yaml:"fetch_remote_list" toml:"fetch_remote_list" xml:"fetchRemoteList" ini:"fetchRemoteList" csv:"FetchRemoteList" json:"fetch_remote_list" yaml:"fetch_remote_list" toml:"fetch_remote_list" xml:"fetchRemoteList" ini:"fetchRemoteList" csv:"FetchRemoteList"`

	// PoolMode
	PoolMode bool `default:"true" flag:"-with-proxy-pool" yaml:"pool_mode" toml:"pool_mode" xml:"poolMode" ini:"poolMode" csv:"PoolMode" json:"pool_mode" yaml:"pool_mode" toml:"pool_mode" xml:"poolMode" ini:"poolMode" csv:"PoolMode"`

	// List
	List []Connect `json:"list" yaml:"list" toml:"list" xml:"list" ini:"list" csv:"list"`

	//-- End
}

// Connect
type Connect struct {

	// Enabled
	Enabled bool `default:"true" flag:"enabled" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled" csv:"Enabled" json:"enabled" yaml:"enabled" toml:"enabled" xml:"enabled" ini:"enabled" csv:"Enabled"`

	// Address
	Address string `required:"true" flag:"address" yaml:"address" toml:"address" xml:"address" ini:"address" csv:"address" json:"address" yaml:"address" toml:"address" xml:"address" ini:"address" csv:"address"`

	//-- End
}
