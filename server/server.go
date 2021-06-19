package server

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
	"xorm.io/xorm"
)

type Config struct {
	Server struct {
		Mysql string `yaml:"mysql"`
	}
	Ywtree struct {
		Host   string `yaml:"host"`
		Secret string `yaml:"secret"`
	} `yaml:"ywtree"`
}

func Run() {
	cfg := &Config{}
	bts, err := ioutil.ReadFile("app.yml")
	if err != nil {
		println("can not read app.yml.please create")
		return
	}
	err = yaml.Unmarshal(bts, cfg)
	if err != nil {
		println("can not format app.yml.")
		return
	}
	xorm.NewEngine("mysql")
}
