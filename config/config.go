package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/daveadams/go-rapture/log"
)

const DefaultAwsRegion = "us-east-1"

type RaptureConfig struct {
	Identifier      string `json:"identifier,omitempty"`
	SessionDuration int64  `json:"session_duration,omitempty"`
	InitMethod      string `json:"init_method,omitempty"`
	DefaultVault    string `json:"default_vault,omitempty"`
	Quiet           bool   `json:"quiet,omitempty"`
}

var config *RaptureConfig

func ConfigFilename() string {
	log.Trace("config: ConfigFilename()")
	return filepath.Join(ConfigDir(), "config.json")
}

func DefaultConfig() *RaptureConfig {
	log.Trace("config: DefaultConfig()")
	return &RaptureConfig{
		Identifier:      os.Getenv("USER"),
		SessionDuration: 3600,
		InitMethod:      "vaulted",
		DefaultVault:    "default",
		Quiet:           false,
	}
}

func LoadConfig() (*RaptureConfig, error) {
	log.Trace("config: LoadConfig()")

	if config != nil {
		return config, nil
	}

	config := DefaultConfig()
	fn := ConfigFilename()
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		// no config file, return defaults
		return config, nil
	} else {
		bytes, err := ioutil.ReadFile(fn)
		if err != nil {
			return config, err
		}
		err = json.Unmarshal(bytes, config)
		if err != nil {
			return config, err
		}
	}
	return config, nil
}

// just return a config even if there's an error
func GetConfig() *RaptureConfig {
	log.Trace("config: GetConfig()")

	if c, err := LoadConfig(); err != nil {
		log.Debugf("Failed to load config: %s", err)
		return DefaultConfig()
	} else {
		return c
	}
}

// return the raw config without any defaults (or an empty config if it's missing)
// return values are config, exists?, error
func RawConfig() (*RaptureConfig, bool, error) {
	log.Trace("config: RawConfig()")

	config := &RaptureConfig{}
	empty := &RaptureConfig{}

	fn := ConfigFilename()
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		// no config file, return an empty config
		log.Debug("Found no config file")
		return empty, false, nil
	} else {
		bytes, err := ioutil.ReadFile(fn)
		if err != nil {
			log.Debugf("Could not read config file: %s", err)
			return empty, true, err
		}
		err = json.Unmarshal(bytes, config)
		if err != nil {
			log.Debugf("Could not parse config: %s", err)
			return empty, true, err
		}
	}

	return config, true, nil
}

func (rc *RaptureConfig) Region() string {
	if rv, ok := os.LookupEnv("AWS_DEFAULT_REGION"); ok {
		return rv
	} else if rv, ok := os.LookupEnv("AWS_REGION"); ok {
		return rv
	} else {
		return DefaultAwsRegion
	}
}
