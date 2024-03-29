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
	v1     *viper.Viper
	appCfg ApplicationCfg
)

type ApplicationCfg struct {
	Gin struct {
		Host string `yaml:"host" json:"host" mapstructure:"host"`
		Port int    `yaml:"port" json:"port" mapstructure:"port"`
		CORS struct {
			AllowOrigins []string `yaml:"allow_origins" json:"allow_origins" mapstructure:"allow_origins"`
			AllowMethods []string `yaml:"allow_methods" json:"allow_methods" mapstructure:"allow_methods"`
		} `yaml:"cors" json:"cors" mapstructure:"cors"`
		Jwt struct {
			SecretKey  string `yaml:"secret_key" json:"secret_key" mapstructure:"secret_key"`
			Expiration int    `yaml:"expiration" json:"expiration" mapstructure:"expiration"`
		} `yaml:"jwt" json:"jwt" mapstructure:"jwt"`
		Paths struct {
			Assets   string `yaml:"assets" json:"assets" mapstructure:"assets"`
			Media    string `yaml:"media" json:"media" mapstructure:"media"`
			Frontend string `yaml:"frontend" json:"frontend" mapstructure:"frontend"`
		} `yaml:"paths" json:"paths" mapstructure:"paths"`
	} `yaml:"gin" json:"gin" mapstructure:"gin"`
	Email struct {
		Address  string `yaml:"address" json:"address" mapstructure:"address"`
		Password string `yaml:"password" json:"password" mapstructure:"password"`
		SMTP     struct {
			Host string `yaml:"host" json:"host" mapstructure:"host"`
			Port int    `yaml:"port" json:"port" mapstructure:"port"`
		} `yaml:"smtp" json:"smtp" mapstructure:"smtp"`
	} `yaml:"email" json:"email" mapstructure:"email"`
	Captcha struct {
		Enabled   bool   `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
		Provider  string `yaml:"provider" json:"provider" mapstructure:"provider"`
		ReCaptcha struct {
			URL       string  `yaml:"url" json:"url" mapstructure:"url"`
			SiteKey   string  `yaml:"site_key" json:"site_key" mapstructure:"site_key"`
			SecretKey string  `yaml:"secret_key" json:"secret_key" mapstructure:"secret_key"`
			Threshold float64 `yaml:"threshold" json:"threshold" mapstructure:"threshold"`
		} `yaml:"re_captcha" json:"re_captcha" mapstructure:"re_captcha"`
		Turnstile struct {
			URL       string `yaml:"url" json:"url" mapstructure:"url"`
			SiteKey   string `yaml:"site_key" json:"site_key" mapstructure:"site_key"`
			SecretKey string `yaml:"secret_key" json:"secret_key" mapstructure:"secret_key"`
		} `yaml:"turnstile" json:"turnstile" mapstructure:"turnstile"`
	} `yaml:"captcha" json:"captcha" mapstructure:"captcha"`
	Db struct {
		Provider string `yaml:"provider" json:"provider" mapstructure:"provider"`
		Postgres struct {
			Host     string `yaml:"host" json:"host" mapstructure:"host"`
			Port     int    `yaml:"port" json:"port" mapstructure:"port"`
			Username string `yaml:"username" json:"username" mapstructure:"username"`
			Password string `yaml:"password" json:"password" mapstructure:"password"`
			Dbname   string `yaml:"dbname" json:"dbname" mapstructure:"dbname"`
			Sslmode  string `yaml:"sslmode" json:"sslmode" mapstructure:"sslmode"`
		} `yaml:"postgres" json:"postgres" mapstructure:"postgres"`
		SQLite3 struct {
			Filename string `yaml:"filename" json:"filename" mapstructure:"filename"`
		} `yaml:"sqlite3" json:"sqlite3" mapstructure:"sqlite3"`
		MySQL struct {
			Host     string `yaml:"host" json:"host" mapstructure:"host"`
			Port     int    `yaml:"port" json:"port" mapstructure:"port"`
			Username string `yaml:"username" json:"username" mapstructure:"username"`
			Password string `yaml:"password" json:"password" mapstructure:"password"`
			Dbname   string `yaml:"dbname" json:"dbname" mapstructure:"dbname"`
		} `yaml:"mysql" json:"mysql" mapstructure:"mysql"`
	} `yaml:"db" json:"db" mapstructure:"db"`
	Container struct {
		Provider string `yaml:"provider" json:"provider" mapstructure:"provider"`
		Nat      struct {
			Type  string `yaml:"type" json:"type" mapstructure:"type"`
			Entry string `yaml:"entry" json:"entry" mapstructure:"entry"`
		}
		Proxy struct {
			Enabled bool   `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
			Type    string `yaml:"type" json:"type" mapstructure:"type"`
			TCP     struct {
				Entry string `yaml:"entry" json:"entry" mapstructure:"entry"`
			}
			WS struct{}
		}
		TrafficCapture struct {
			Enabled bool   `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
			Path    string `yaml:"path" json:"path" mapstructure:"path"`
		} `yaml:"traffic_capture" json:"traffic_capture" mapstructure:"traffic_capture"`
		Docker struct {
			URI   string `yaml:"uri" json:"uri" mapstructure:"uri"`
			Entry string `yaml:"entry" json:"entry" mapstructure:"entry"`
		} `yaml:"docker" json:"docker" mapstructure:"docker"`
		K8s struct {
			Namespace string `yaml:"namespace" json:"namespace" mapstructure:"namespace"`
			Path      struct {
				Config string `yaml:"config" json:"config" mapstructure:"config"`
			} `yaml:"path" json:"path" mapstructure:"path"`
		} `yaml:"k8s" json:"k8s" mapstructure:"k8s"`
	} `yaml:"container" json:"container" mapstructure:"container"`
}

func AppCfg() *ApplicationCfg {
	return &appCfg
}

func InitApplicationCfg() {
	v1 = viper.New()
	configFile := path.Join("application.json")
	v1.SetConfigType("json")
	v1.SetConfigFile(configFile)
	if _, err := os.Stat(configFile); err != nil {
		zap.L().Warn("No configuration file found, default configuration file will be created.")

		// Read default configuration from embed
		defaultConfig, _err := embed.FS.Open("configs/application.json")
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
			zap.L().Fatal("Unable to create default configuration file.")
		}
		zap.L().Info("The default configuration file has been generated.")
	}

	if err := v1.ReadInConfig(); err != nil {
		zap.L().Fatal("Unable to read configuration file.")
		return
	}

	if err := v1.Unmarshal(&appCfg); err != nil {
		zap.L().Fatal("Unable to parse configuration file to structure.")
	}

	Mkdirs()
}

func Mkdirs() {
	if AppCfg().Container.TrafficCapture.Enabled {
		if _, err := os.Stat(AppCfg().Container.TrafficCapture.Path); err != nil {
			if _err := os.MkdirAll(AppCfg().Container.TrafficCapture.Path, os.ModePerm); _err != nil {
				zap.L().Fatal("Unable to create directory for traffic capture.")
			}
		}
	}

	if _, err := os.Stat(AppCfg().Gin.Paths.Assets); err != nil {
		if _err := os.MkdirAll(AppCfg().Gin.Paths.Assets, os.ModePerm); _err != nil {
			zap.L().Fatal("Unable to create directory for assets.")
		}
	}

	if _, err := os.Stat(AppCfg().Gin.Paths.Media); err != nil {
		if _err := os.MkdirAll(AppCfg().Gin.Paths.Media, os.ModePerm); _err != nil {
			zap.L().Fatal("Unable to create directory for media.")
		}
	}

	if _, err := os.Stat(AppCfg().Gin.Paths.Frontend); err != nil {
		if _err := os.MkdirAll(AppCfg().Gin.Paths.Frontend, os.ModePerm); _err != nil {
			zap.L().Fatal("Unable to create directory for frontend.")
		}
	}
}

func (a *ApplicationCfg) Save() (err error) {
	val := reflect.ValueOf(appCfg)
	typeOfCfg := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		v1.Set(typeOfCfg.Field(i).Tag.Get("mapstructure"), field.Interface())
	}
	err = v1.WriteConfig()
	return err
}
