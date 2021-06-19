package comm

import "xorm.io/xorm"

var (
	Cfg = &Config{}
	Db  *xorm.Engine
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
