package prop

import "log"

type DatabaseProp struct {
	*Database `yaml:"database" json:"database"`
}

type Database struct {
	Postgresql *DataBaseInfo `yaml:"postgresql" json:"postgresql"`
	Mysql      *DataBaseInfo `yaml:"mysql" json:"mysql"`
}

type DataBaseInfo struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Name     string `yaml:"database" json:"database"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
}

func (d *DatabaseProp) SetDefault() {
	if d.Postgresql.Host == "" {
		log.Fatalf("Yml file parameter 'database.postgresql.host' is empty.")
	}
	if d.Postgresql.Port == 0 {
		d.Postgresql.Port = 5432
	}
}
