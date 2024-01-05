package config

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/erbieio/web2-bridge/utils/logger"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

const (
	LOGOUT_FILE   = "file"
	LOGOUT_STDOUT = "stdout"
)

// one DB one instance
type RedisConfig struct {
	UseCA        bool
	IP           string
	ConnPort     int64
	SSHPort      int64
	SSHAccount   string
	SSHKey       string
	Name         string
	Password     string
	Host         string
	DB           int64
	MinIdleConns int64
}

// one database one instance
type MysqlConfig struct {
	UseCA        bool
	IP           string
	ConnPort     int64
	SSHPort      int64
	SSHAccount   string
	SSHKey       string
	Account      string
	Password     string
	SqlName      string
	MaxOpenConns int64
	MaxIdleConns int64
	MaxLifetime  int64
}
type ServerConfig struct {
	RunMode             string
	HttpPort            string
	ReadTimeout         int64
	WriteTimeout        int64
	TLSCAFile           string
	TLSCAKey            string
	JWTExpireTimeMinute int64
	JWTSecret           string
	LogOut              string
}

type TwitterConfig struct {
	OauthToken       string
	OauthTokenSecret string
	ApiKey           string
	ApiKeySecret     string
	Bearer           string
}

type DiscordConfig struct {
	BotToken string
}

type IpfsConfig struct {
	Api         string
	HttpGateway string
}

type ChainConfig struct {
	Rpc          string
	NftAdminPriv string
	MaxMint      int
}

type GradioConfig struct {
	Url string
}

type ComfyuiConfig struct {
	Host string
	Port int
}

// struct decode must has tag
type Config struct {
	ServerConf  ServerConfig  `toml:"ServerConfig" mapstructure:"ServerConfig"`
	MysqlConf   MysqlConfig   `toml:"MysqlConfig" mapstructure:"MysqlConfig"`
	RedisConf   RedisConfig   `toml:"RedisConfig" mapstructure:"RedisConfig"`
	TwitterConf TwitterConfig `toml:"TwitterConfig" mapstructure:"TwitterConfig"`
	DiscordConf DiscordConfig `toml:"DiscordConfig" mapstructure:"DiscordConfig"`
	IpfsConf    IpfsConfig    `toml:"IpfsConfig" mapstructure:"IpfsConfig"`
	ChainConf   ChainConfig   `toml:"ChainConfig" mapstructure:"ChainConfig"`
	GradioConf  GradioConfig  `toml:"GradioConfig" mapstructure:"GradioConfig"`
	ComfyuiConf ComfyuiConfig `toml:"ComfyuiConfig" mapstructure:"ComfyuiConfig"`
}

var (
	configMutex     = sync.RWMutex{}
	configPath      = ""
	configFileAbs   = ""
	config          Config
	configViper     *viper.Viper
	configFlyChange []chan bool
)

func RegistConfChange(c chan bool) {
	configFlyChange = append(configFlyChange, c)
}

func notifyConfChange() {
	for i := 0; i < len(configFlyChange); i++ {
		configFlyChange[i] <- true
	}
}

func watchConfig(c *viper.Viper) error {
	c.WatchConfig()
	c.OnConfigChange(func(e fsnotify.Event) {
		logger.Logrus.WithFields(logrus.Fields{"change": e.String()}).Info("config change and reload it")
		reloadConfig(c)
		notifyConfChange()
	})
	return nil
}

func LoadConf(configFilePath string) error {
	config = Config{}
	configMutex.Lock()
	defer configMutex.Unlock()

	configViper = viper.New()
	configViper.SetConfigName("config")
	configViper.AddConfigPath(configFilePath) //endwith "/"
	configViper.SetConfigType("yaml")

	if err := configViper.ReadInConfig(); err != nil {
		return err
	}
	if err := configViper.Unmarshal(&config); err != nil {
		return err
	}

	s, _ := json.MarshalIndent(config, "", "\t")
	fmt.Printf("Load config: %s", s)

	if err := watchConfig(configViper); err != nil {
		return err
	}
	return nil
}

func reloadConfig(c *viper.Viper) {
	configMutex.Lock()
	defer configMutex.Unlock()

	if err := c.ReadInConfig(); err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("config ReLoad failed")
	}

	if err := configViper.Unmarshal(&config); err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("unmarshal config failed")
	}

	logger.Logrus.WithFields(logrus.Fields{"config": config}).Info("Config ReLoad Success")
}

func GetServerConfig() ServerConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.ServerConf
}

func GetMysqlConfig() MysqlConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.MysqlConf
}

func GetRedisConfig() RedisConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.RedisConf
}
func GetTwitterConfig() TwitterConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.TwitterConf
}

func GetDiscordConfig() DiscordConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.DiscordConf
}

func GetIpfsConfig() IpfsConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.IpfsConf
}

func GetChainConfig() ChainConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.ChainConf
}

func GetGradioConfig() GradioConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.GradioConf
}

func GetComfyuiConfig() ComfyuiConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.ComfyuiConf
}

// check if logout equal file
func (c ServerConfig) LogOutFile() bool {
	return c.LogOut == LOGOUT_FILE
}

// check if logout equal stdout
func (c ServerConfig) LogOutStdout() bool {
	return c.LogOut == LOGOUT_STDOUT
}
