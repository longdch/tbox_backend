package config

import (
	"bytes"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"log"
	"strings"
)

var defaultConfig = []byte(`
base:
  environment: Local
  port: 8080
mysql:
  address: tbox_database:3306
  protocol: tcp
  database: tbox
  username: dchlong
  password: dchlong
  multi_statements: true
  allow_native_passwords: true
  parse_time: true
phone_number_rate_limit:
  limit: 10
  burst: 10
otp:
  expired_time: 60
  resend_waiting_time: 30
  size: 6
sms_service:
  url: https://5db83e44177b350014ac77c6.mockapi.io/v1/sms
swagger:
  url: http://localhost:8080/swagger/doc.json
token:
  secret_key: 5OQ3ldRoOlkFg5PavqYXlWTZ88gc1DPE
`)

type Config struct {
	Base
	MySQL                MySQL                `yaml:"mysql" mapstructure:"mysql"`
	PhoneNumberRateLimit PhoneNumberRateLimit `yaml:"phone_number_rate_limit" mapstructure:"phone_number_rate_limit"`
	Otp                  Otp                  `yaml:"otp" mapstructure:"otp"`
	SmsService           SmsService           `yaml:"sms_service" mapstructure:"sms_service"`
	Token                Token                `yaml:"token" mapstructure:"token"`
	Swagger              Swagger              `yaml:"swagger" mapstructure:"swagger"`
}

type MySQL struct {
	Username             string `yaml:"username" mapstructure:"username"`
	Password             string `yaml:"password" mapstructure:"password"`
	Address              string `yaml:"address" mapstructure:"address"`
	Protocol             string `yaml:"protocol" mapstructure:"protocol"`
	Database             string `yaml:"database" mapstructure:"database"`
	MultiStatements      bool   `yaml:"multi_statements" mapstructure:"multi_statements"`
	AllowNativePasswords bool   `yaml:"allow_native_passwords" mapstructure:"allow_native_passwords"`
	ParseTime            bool   `yaml:"parse_time" mapstructure:"parse_time"`
}

type PhoneNumberRateLimit struct {
	Limit float64 `yaml:"limit" mapstructure:"limit"`
	Burst int     `yaml:"burst" mapstructure:"burst"`
}

type Otp struct {
	ExpiredTime       int `yaml:"expired_time" mapstructure:"expired_time"`
	ResendWaitingTime int `yaml:"resend_waiting_time" mapstructure:"resend_waiting_time"`
	Size              int `yaml:"size" mapstructure:"size"`
}

type SmsService struct {
	Url string `yaml:"url" mapstructure:"url"`
}

type Swagger struct {
	Url string `yaml:"url" mapstructure:"url"`
}

type Token struct {
	SecretKey string `yaml:"secret_key" mapstructure:"secret_key"`
}

// FormatDSN returns MySQL DSN from settings.
func (m *MySQL) FormatDSN() string {
	um := &mysql.Config{
		User:                 m.Username,
		Passwd:               m.Password,
		Addr:                 m.Address,
		Net:                  m.Protocol,
		DBName:               m.Database,
		MultiStatements:      m.MultiStatements,
		AllowNativePasswords: m.AllowNativePasswords,
		ParseTime:            m.ParseTime,
	}

	return um.FormatDSN()
}

func (cfg *Config) MySqlDSN() string {
	return cfg.MySQL.FormatDSN()
}

type Tracing struct {
	Enabled bool
}

type Base struct {
	Environment string `yaml:"environment"`
	Port        int    `yaml:"port"`
}

func Load() Config {
	var cfg = Config{}
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	if err != nil {
		log.Fatalf("Failed to read config %v", err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("Failed to unmarshal config %v", err)
	}

	log.Println("Config loaded: ", cfg)
	return cfg
}
