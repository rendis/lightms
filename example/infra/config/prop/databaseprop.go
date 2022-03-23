package prop

import "sync"

var databaseProp *DatabaseProp
var dataBaseOnce sync.Once

func GetDatabaseProp() *DatabaseProp {
	dataBaseOnce.Do(func() {
		databaseProp = &DatabaseProp{}
	})
	return databaseProp
}

type DatabaseProp struct {
	Database `yaml:"database"`
}

type Database struct {
	Postgresql DataBaseInfo `yaml:"postgresql"`
	Mysql      DataBaseInfo `yaml:"mysql"`
}

type DataBaseInfo struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
