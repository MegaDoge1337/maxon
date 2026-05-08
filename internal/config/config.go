package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Mode       string     `mapstructure:"mode"`
	Env        string     `mapstructure:"env"`
	HTTPServer HTTPServer `mapstructure:"http_server"`
	Database   Database   `mapstructure:"database"`
	Auth       Auth       `mapstructure:"auth"`
}

type HTTPServer struct {
	Address     string        `mapstructure:"address"`
	Timeout     time.Duration `mapstructure:"timeout"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}

type Database struct {
	Driver          string `mapstructure:"db_driver"`
	Host            string `mapstructure:"db_host"`
	Port            int    `mapstructure:"db_port"`
	User            string `mapstructure:"db_user"`
	Password        string `mapstructure:"db_password"`
	Name            string `mapstructure:"db_name"`
	SslMode         string `mapstructure:"db_sslmode"`
	MigrationsTable string `mapstructure:"db_migrations_table"`
}

type Auth struct {
	JwtSecret string `mapstructure:"jwt_secret"`
}

func LoadConfig(cfg interface{}) error {
	viper.SetConfigType("yaml")
	viper.SetConfigFile("./config/config.yml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	env := viper.New()
	env.SetConfigType("env")
	env.SetConfigFile(viper.GetString("env"))
	if err := env.ReadInConfig(); err != nil {
		return err
	}

	nested := map[string]interface{}{
		"database": map[string]interface{}{
			"db_driver":           env.GetString("db_driver"),
			"db_host":             env.GetString("db_host"),
			"db_port":             env.GetInt("db_port"),
			"db_user":             env.GetString("db_user"),
			"db_password":         env.GetString("db_password"),
			"db_name":             env.GetString("db_name"),
			"db_sslmode":          env.GetString("db_sslmode"),
			"db_migrations_table": env.GetString("db_migrations_table"),
		},
		"auth": map[string]interface{}{
			"jwt_secret": env.GetString("jwt_secret"),
		},
	}

	if err := viper.MergeConfigMap(nested); err != nil {
		return err
	}

	return viper.Unmarshal(cfg)
}
