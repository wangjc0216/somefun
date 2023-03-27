package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	OpenaiKey                 string `json:"openai_key"`
	OpenaiMaxConcurrencyLimit int    `json:"openai_max_concurrency_limit"`
}

var c Config

func LoadConfig(filepath string) *Config {
	bs, err := os.ReadFile(filepath)
	if err != nil {
		panic(fmt.Errorf("LoadConfig error:%w", err))
	}

	err = json.Unmarshal(bs, &c)
	if err != nil {
		panic(fmt.Errorf("json unmarshal error:%w", err))
	}
	//todo check required field and  Optional field
	return &c
}
func GetConfig() *Config {
	return &c
}
