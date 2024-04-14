package config

import (
	"github.com/go-ini/ini"
	"log"
)

type App struct {
	APIHOST                string
	APIPORT                string
	JWTSECRET              string
	LogFilePath            string
	LogFile                string
	AccessTokenExpiredTime int
}

var AppSetting = &App{}

type Database struct {
	DBURL string
}

var DatabaseSetting = &Database{}
var cfg *ini.File

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}

type Config struct {
	App      App
	Database Database
}

func ProvideConfig() *Config {
	var err error
	//"./conf/app.ini"
	cfg, err = ini.Load("./conf/app.env")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}
	mapTo("app", AppSetting)
	mapTo("database", DatabaseSetting)
	return &Config{
		App:      *AppSetting,
		Database: *DatabaseSetting,
	}
}
