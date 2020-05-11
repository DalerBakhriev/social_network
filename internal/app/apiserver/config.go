package apiserver

// Config contains apiserver configuration setting
type Config struct {
	BindAddr   string `toml:"bind_addr"`
	LogLevel   string `toml:"log_level"`
	SessionKey string `toml:"session_key"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}
