package config

import (
	"github.com/go-ini/ini"
	"log"
	"os"
	"time"
)

type server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type mongoDb struct {
	URI    string
	DBName string
}

var MongoDB = &mongoDb{}
var Server = &server{}
var cfg *ini.File

func init() {
	var err error
	cfg, err = ini.Load(os.Getenv("GOPATH") + "/.env")
	if err != nil {
		log.Fatalf("Fail to parse '.env': %v", err)
	}

	mapTo("server", Server)
	mapTo("mongodb", MongoDB)

	Server.ReadTimeout = Server.ReadTimeout * time.Second
	Server.WriteTimeout = Server.ReadTimeout * time.Second
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf(".env -> config mapping error: %v", err)
	}
}
