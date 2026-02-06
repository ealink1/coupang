package config

import (
	"flag"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var configOnce sync.Once

// 定义 -f 命令行标志，默认值为 "config.yaml"
var cfg *Config

var configFile = flag.String("f", "config.yaml", "Path to the YAML config file")

func GetCfg() *Config {
	if cfg == nil {
		configOnce.Do(func() {
			// 读取配置文件
			data, err := os.ReadFile(*configFile)
			if err != nil {
				log.Fatalf("Failed to read config file %s: %v", *configFile, err)
			}

			// 解析 YAML
			var config Config
			if err = yaml.Unmarshal(data, &config); err != nil {
				log.Fatalf("Failed to parse YAML config: %v", err)
			}
			cfg = &config
		})
	}
	return cfg
}

type Config struct {
	Coupang CoupangConfig `yaml:"coupang"`
}

type CoupangConfig struct {
	VendorId  string `yaml:"vendor_id"`
	ApiKey    string `yaml:"api_key"`
	SecretKey string `yaml:"secret_key"`
}
