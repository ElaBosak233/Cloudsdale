package config

import (
	"github.com/elabosak233/cloudsdale/embed"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"os"
	"path"
	"reflect"
)

var (
	v2     *viper.Viper
	pltCfg PlatformCfg
)

type PlatformCfg struct {
	Site struct {
		Title       string `yaml:"title" json:"title" mapstructure:"title"`
		Description string `yaml:"description" json:"description" mapstructure:"description"`
	} `yaml:"site" json:"site" mapstructure:"site"`
	Container struct {
		ParallelLimit int `yaml:"parallel_limit" json:"parallel_limit" mapstructure:"parallel_limit"`
		RequestLimit  int `yaml:"request_limit" json:"request_limit" mapstructure:"request_limit"`
	} `yaml:"container" json:"container" mapstructure:"container"`
	User struct {
		AllowRegistration bool `yaml:"allow_registration" json:"allow_registration" mapstructure:"allow_registration"`
	} `yaml:"user" json:"user" mapstructure:"user"`
}

func PltCfg() *PlatformCfg {
	return &pltCfg
}

func InitPlatformCfg() {
	v2 = viper.New()
	configFile := path.Join("platform.json")
	v2.SetConfigType("json")
	v2.SetConfigFile(configFile)
	if _, err := os.Stat(configFile); err != nil {
		zap.L().Warn("No configuration file found, default configuration file will be created.")

		// Read default configuration from embed
		defaultConfig, _err := embed.FS.Open("configs/platform.json")
		if _err != nil {
			zap.L().Error("Unable to read default configuration file.")
			return
		}
		defer func(defaultConfig fs.File) {
			_ = defaultConfig.Close()
		}(defaultConfig)

		// Create config file in current directory
		dstConfig, _err := os.Create(configFile)
		defer func(dstConfig *os.File) {
			_ = dstConfig.Close()
		}(dstConfig)

		if _, _err = io.Copy(dstConfig, defaultConfig); _err != nil {
			zap.L().Error("Unable to create default configuration file.")
		}
		zap.L().Info("The default configuration file has been generated.")
	}

	if err := v2.ReadInConfig(); err != nil {
		zap.L().Error("Unable to read configuration file.")
		return
	}

	if err := v2.Unmarshal(&pltCfg); err != nil {
		zap.L().Error("Unable to parse configuration file to structure.")
	}
}

func (p *PlatformCfg) Save() (err error) {
	val := reflect.ValueOf(pltCfg)
	typeOfCfg := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		v2.Set(typeOfCfg.Field(i).Tag.Get("mapstructure"), field.Interface())
	}
	err = v2.WriteConfig()
	return err
}
