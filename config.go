package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	DbHost                  string `json:"db_host"`
	DbPort                  int    `json:"db_port"`
	DbUser                  string `json:"db_user"`
	DbPassword              string `json:"db_password"`
	DbName                  string `json:"db_name"`
	DbConnectionMaxLifeTime int    `json:"db_connection_max_life_time"`
	DbMaxOpenConnections    int    `json:"db_max_open_connections"`
	DbMaxIdleConnections    int    `json:"db_max_idle_connections"`

	WebServerHost string `json:"webserver_host"`
	WebServerPort int    `json:"webserver_port"`

	SessionCookieName string `json:"session_cookie_name"`
	SessionMaxAge     int    `json:"session_max_age"`
	SessionStoreKey   string `json:"session_store_key"`
	CookieStoreKey    string `json:"cookie_store_key"`
}

func NewConfig(config_path string) (*Config, error) {
	cfg := new(Config)
	return cfg, cfg.init(config_path)
}

func (cfg *Config) init(config_path string) (err error) {
	text, err := ioutil.ReadFile(config_path)
	if err == nil {
		err = json.Unmarshal(text, cfg)
		if err != nil {
			return err
		}
	}

	return err
}
