package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Watcher struct {
	PATH      string
	TIMEOUT   int
	STABLE    int
	RUNCONFIG K8RunConfig
}

type K8RunConfig struct {
	NAMESPACE   string
	COMMAND     []string
	IMAGE       string
	ENVIRONMENT map[string]string
}

type PersistentStore struct {
	REDISHOST string
	REDISDB   int
	REDISPASS string
}

type ConfigFile struct {
	WATCHERS []Watcher
	REDIS    PersistentStore
}

func LoadConfig(filename string) (*ConfigFile, error) {
	configData := ConfigFile{}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	marshalErr := yaml.Unmarshal([]byte(data), &configData)
	if marshalErr != nil {
		return nil, marshalErr
	}

	return &configData, nil
}

func WatcherFor(path string, config *ConfigFile) (*Watcher, error) {
	for _, v := range config.WATCHERS {
		if v.PATH == path {
			return &v, nil
		}
	}
	return nil, nil
}
