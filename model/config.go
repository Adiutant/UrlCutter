package model

import (
	"encoding/json"
	"fmt"
	"os"
)

type UrlServerConfig struct {
	CacheSize int `json:"cache_size"`
}

func (c *UrlServerConfig) ReadConfig(postfix string) error {
	file, err := os.ReadFile(fmt.Sprintf("./config/config_%s.json", postfix))
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, c)
	if err != nil {
		return err
	}
	return nil
}
