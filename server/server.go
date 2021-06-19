package server

import (
	"errors"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/yggworldtree/cpu_report/comm"
	"github.com/yggworldtree/go-sdk/ywtree"
	"gopkg.in/yaml.v3"
	"xorm.io/xorm"
)

var (
	Mgr    *Manager
	YwtEgn *ywtree.Engine
)

func Run() {
	hbtp.Debug = true
	bts, err := ioutil.ReadFile("app.yml")
	if err != nil {
		println("can not read app.yml.please create")
		return
	}
	err = yaml.Unmarshal(bts, comm.Cfg)
	if err != nil {
		println("can not format app.yml.")
		return
	}
	err = InitXorm(comm.Cfg.Server.Mysql, &comm.Db)
	if err != nil {
		println("init mysql err:" + err.Error())
		return
	}
	defer comm.Db.Close()
	Mgr = NewManager()
	YwtEgn = ywtree.NewEngine(Mgr, &ywtree.Config{
		Host:   comm.Cfg.Ywtree.Host,
		Secret: comm.Cfg.Ywtree.Secret,
		Org:    "mgr",
		Name:   "cpu-report",
	})
	err = YwtEgn.Run()
	if err != nil {
		println("ywtree err:" + err.Error())
		return
	}
}
func InitXorm(ul string, pdb **xorm.Engine) error {
	if ul == "" {
		return errors.New("url blank")
	}
	db, err := xorm.NewEngine("mysql", ul)
	if err != nil {
		return err
	}
	*pdb = db
	// *pdb = gocloud.NewDBHelper(db)
	return nil
}
