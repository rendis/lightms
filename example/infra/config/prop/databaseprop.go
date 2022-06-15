package prop

import "log"

type DatabaseProp struct {
	*Database `yaml:"database"`
}

type Database struct {
	Postgresql *DataBaseInfo `yaml:"postgresql"`
	Mysql      *DataBaseInfo `yaml:"mysql"`
}

type DataBaseInfo struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func (d *DatabaseProp) SetDefault() {
	if d.Postgresql.Host == "" {
		log.Fatalf("Yml file parameter 'database.postgresql.host' is empty.")
	}
	if d.Postgresql.Port == 0 {
		d.Postgresql.Port = 5432
	}
}
