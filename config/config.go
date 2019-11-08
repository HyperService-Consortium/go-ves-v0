package config

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
)

var (
	defaultDatabaseConfig = DatabaseConfig{
		Engine:         "xorm",
		Type:           "mysql",
		UserName:       "ves",
		Password:       "123456",
		ConnectionType: "tcp",
		RemoteHost:     "127.0.0.1:3306",
		BaseName:       "ves",
		Encoding:       "utf8",
	}
	defaultKVDBConfig = KVDBConfig{
		Type: "leveldb",
		Path: "./data/",
	}
	defaultServerConfig = ServerConfig{
		Port:              ":23351",
		CentralVesAddress: "127.0.0.1:23352",
	}

	defaultConfig = &Configuration{
		defaultDatabaseConfig,
		defaultKVDBConfig,
		defaultServerConfig,
	}
	cfg         *Configuration
	CfgContext  string //= flag.String("config", "./ves-server-config.toml", "configurate")
	cfgLock     sync.RWMutex
	parseConfig sync.Once
)

// Configuration defines env-vari of ves-server
type Configuration struct {
	DatabaseConfig DatabaseConfig `toml:"database"`
	KVDBConfig     KVDBConfig     `toml:"kvdb"`
	ServerConfig   ServerConfig   `toml:"server"`
}

// DatabaseConfig defines configuration of database.*
type DatabaseConfig struct {
	Engine         string `toml:"engine"`
	Type           string `toml:"type"`
	UserName       string `toml:"user_name"`
	Password       string `toml:"password"`
	ConnectionType string `toml:"connection_type"`
	RemoteHost     string `toml:"host"`
	BaseName       string `toml:"base_name"`
	Encoding       string `toml:"encoding"`
}

// KVDBConfig defines configuration of kvdb.*
type KVDBConfig struct {
	Type string `toml:"type"`
	Path string `toml:"path"`
}

// ServerConfig defines configuration of server.*
type ServerConfig struct {
	Port              string `toml:"port"`
	CentralVesAddress string `toml:"central_ves_address"`
}

// Config return a singleton of configuration struct
func Config() *Configuration {
	parseConfig.Do(func() { ReloadConfiguration() })
	cfgLock.RLock()
	defer cfgLock.RUnlock()
	return cfg
}

// ResetPath will reset path and reload configuration from the path
func ResetPath(path string) {
	CfgContext = path
	ReloadConfiguration()
}

// ReloadConfiguration flushes the singleton to newest status
func ReloadConfiguration() error {
	filePath, err := filepath.Abs(CfgContext)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("reseting: %s\n", filePath)
	config := new(Configuration)
	*config = *defaultConfig
	if _, err := toml.DecodeFile(filePath, config); err != nil {
		panic(fmt.Errorf("err: %v, %v", filePath, err))
	}
	cfgLock.Lock()
	defer cfgLock.Unlock()
	cfg = config
	return nil
}
