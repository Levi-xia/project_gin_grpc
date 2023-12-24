package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

var GlobalConf = &Configuration{}

type Config struct {
	viper *viper.Viper
}

type Configuration struct {
	Zap   ZapConfig   `mapstructure:"zap" json:"zap" yaml:"zap"`
	MySql MySqlConfig `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Jwt   JwtConfig   `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
}

type ZapConfig struct {
	DebugFileName string `mapstructure:"debugFileName" json:"debugFileName" yaml:"debugFileName"`
	InfoFileName  string `mapstructure:"infoFileName" json:"infoFileName" yaml:"infoFileName"`
	WarnFileName  string `mapstructure:"warnFileName" json:"warnFileName" yaml:"warnFileName"`
	ErrorFileName string `mapstructure:"errorFileName" json:"errorFileName" yaml:"errorFileName"`
	MaxSize       int    `mapstructure:"maxSize" json:"maxSize" yaml:"maxSize"`
	MaxAge        int    `mapstructure:"maxAge" json:"maxAge" yaml:"maxAge"`
	MaxBackups    int    `mapstructure:"maxBackups" json:"maxBackups" yaml:"maxBackups"`
}

type MySqlConfig struct {
	Driver              string `mapstructure:"driver" json:"driver" yaml:"driver"`
	Host                string `mapstructure:"host" json:"host" yaml:"host"`
	Port                int    `mapstructure:"port" json:"port" yaml:"port"`
	Database            string `mapstructure:"database" json:"database" yaml:"database"`
	Username            string `mapstructure:"username" json:"username" yaml:"username"`
	Password            string `mapstructure:"password" json:"password" yaml:"password"`
	Charset             string `mapstructure:"charset" json:"charset" yaml:"charset"`
	MaxIdleConns        int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns        int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	LogMode             string `mapstructure:"log_mode" json:"log_mode" yaml:"log_mode"`
	EnableFileLogWriter bool   `mapstructure:"enable_file_log_writer" json:"enable_file_log_writer" yaml:"enable_file_log_writer"`
	LogFilename         string `mapstructure:"log_filename" json:"log_filename" yaml:"log_filename"`
}

type JwtConfig struct {
	Secret string `mapstructure:"secret" json:"secret" yaml:"secret"`
	JwtTtl int64  `mapstructure:"jwt_ttl" json:"jwt_ttl" yaml:"jwt_ttl"` // token 有效期（秒）
}

func InitConfig() {
	config := &Config{
		viper: viper.New(),
	}
	config.viper.SetConfigName("app-dev")
	config.viper.SetConfigType("yaml")
	config.viper.AddConfigPath("/Users/levi/go/src/project-gin-grpc/project-common/config")

	if err := config.viper.ReadInConfig(); err != nil {
		log.Fatalln("Fatal error config file: ", err)
	}
	config.viper.WatchConfig()
	config.viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("config file changed:", in.Name)
		if err := config.viper.Unmarshal(&GlobalConf); err != nil {
			log.Println("Unmarshal config failed, err:", err)
		}
	})
	if err := config.viper.Unmarshal(GlobalConf); err != nil {
		log.Println("Unmarshal config failed, err:", err)
	}
}
