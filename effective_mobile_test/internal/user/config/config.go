package config

type Config struct {
	Database   Database   `yaml:"database"`
	HttpServer HttpServer `yaml:"httpServer"`
}

type HttpServer struct {
	Port int `yaml:"port"`
}

type Database struct {
	Main    DbNone
	Replica DbNone
}

type DbNone struct {
	Host     string `yaml:"Host"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	DbName   string `yaml:"DbName"`
	SslMode  string `yaml:"SslMode"`
}
