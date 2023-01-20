package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-ini/ini"
)

type Server struct {
	Address string
	Port    string
}

type Database struct {
	Type     string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

type Logging struct {
	Level    string
	Format   string
	Filename string
}

type JWT struct {
	Secret    string
	ExpiresIn string
}

type Payment struct {
	Gateway              string
	StripeSecretKey      string
	StripePublishableKey string
}

type Shipping struct {
	Provider           string
	FedexAPIKey        string
	FedexAPIPassword   string
	FedexAccountNumber string
	FedexMeterNumber   string
}

type Config struct {
	Server   Server
	Database Database
	Logging  Logging
	JWT      JWT
	Payment  Payment
	Shipping Shipping
}

func (c *Config) Read(file string) error {
	ini, err := ini.Load(file)
	if err != nil {
		return err
	}

	c.Server.Address = ini.Section("server").Key("address").String()
	c.Server.Port = ini.Section("server").Key("port").String()

	c.Database.Type = ini.Section("database").Key("type").String()
	c.Database.Host = ini.Section("database").Key("host").String()
	c.Database.Port = ini.Section("database").Key("port").String()
	c.Database.Name = ini.Section("database").Key("name").String()
	c.Database.User = ini.Section("database").Key("user").String()
	c.Database.Password = ini.Section("database").Key("password").String()
	c.Database.SSLMode = ini.Section("database").Key("ssl_mode").String()

	c.Logging.Level = ini.Section("logging").Key("level").String()
	c.Logging.Format = ini.Section("logging").Key("format").String()
	c.Logging.Filename = ini.Section("logging").Key("filename").String()

	c.JWT.Secret = ini.Section("jwt").Key("secret").String()
	c.JWT.ExpiresIn = ini.Section("jwt").Key("expires_in").String()

	c.Payment.Gateway = ini.Section("payment").Key("gateway").String()
	c.Payment.StripeSecretKey = ini.Section("payment").Key("stripe_secret_key").String()
	c.Payment.StripePublishableKey = ini.Section("payment").Key("stripe_publishable_key").String()

	c.Shipping.Provider = ini.Section("shipping").Key("provider").String()
	c.Shipping.FedexAPIKey = ini.Section("shipping").Key("fedex_api_key").String()
	c.Shipping.FedexAPIPassword = ini.Section("shipping").Key("fedex_api_password").String()
	c.Shipping.FedexAccountNumber = ini.Section("shipping").Key("fedex_account_number").String()
	c.Shipping.FedexMeterNumber = ini.Section("shipping").Key("fedex_meter_number").String()
	return nil
}

func (c Config) DSN() string {
	if strings.ToLower(c.Database.Type) == "mysql" {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name)
	}
	return ""
}

func (c Config) GetJWTSecret() string {
	return c.JWT.Secret
}

func (c Config) GetJWTExpiresIn() (time.Duration, error) {
	timeDur, err := time.ParseDuration(c.JWT.ExpiresIn)
	if err != nil {
		return 0, err
	}
	return timeDur, nil
}

func SetConfig(file string) (*Config, error) {
	c := &Config{}
	if err := c.Read(file); err != nil {
		return nil, err
	}
	return c, nil
}

func GetTestConfig() *Config {
	return &Config{
		Database: Database{
			Host:     "localhost",
			Port:     "3306",
			User:     "root",
			Password: "root",
			Name:     "ecomdb",
		},
	}
}
