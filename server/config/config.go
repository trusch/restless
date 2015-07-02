package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	JSBase    string     `json:"jsBase"`
	Assets    string     `json:"assets"`
	Address   string     `json:"address"`
	DB        string     `json:"db"`
	Endpoints []Endpoint `json:"endpoints"`
	CertFile  string     `json:"certfile"`
	KeyFile   string     `json:"keyfile"`
}

type Endpoint struct {
	Url   string `json:"url"`
	Model string `json:"model"`
}

func Load(path string) *Config {
	cfg := new(Config)
	file, e := ioutil.ReadFile(path)
	if e != nil {
		log.Fatal("File error:", e)
	}
	e = json.Unmarshal(file, &cfg)
	if e != nil {
		log.Fatal("File error:", e)
	}
	return cfg
}
