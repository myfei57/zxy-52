package store

import (
	"path/filepath"
	"time"
)

type PlantConfig struct {
	PlantName string
	Region    string
	Timezone  string
	UpdatedAt string
}

type ConfigStore struct {
	store *Store
}

func NewConfigStore(store *Store) *ConfigStore {
	return &ConfigStore{store: store}
}

func (c *ConfigStore) Path() string {
	return filepath.Join(c.store.Dir("config"), "plant.json")
}

func (c *ConfigStore) Save(config PlantConfig) error {
	config.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return writeJSON(c.Path(), config)
}

func (c *ConfigStore) Load() (PlantConfig, error) {
	var config PlantConfig
	err := readJSON(c.Path(), &config)
	return config, err
}

func (c *ConfigStore) Default() PlantConfig {
	return PlantConfig{
		PlantName: "WasteGen DCS",
		Region:    "east",
		Timezone:  "Asia/Shanghai",
	}
}
