package server

import (
	"errors"
	"time"

	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/yggworldtree/cpu_report/comm"
	"github.com/yggworldtree/cpu_report/model"
	"github.com/yggworldtree/cpu_report/service"
)

type cmds struct {
}

func (c *cmds) AuthFun() hbtp.AuthFun {
	return nil
}

func (cmds) SetWarnInterval(c *hbtp.Context, m *hbtp.Map) {
	val, _ := m.GetInt("value")
	if val <= 0 {
		c.ResString(hbtp.ResStatusErr, "param err")
		return
	}

	bts := hbtp.BigIntToByte(val, 8)
	err := service.SetParam("warn-interval", bts)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, "set param err:"+err.Error())
		return
	}
	Mgr.warnInterval = time.Duration(val)
	Mgr.clearTmr()
	c.ResString(hbtp.ResStatusOk, "ok")
}

func (cmds) SetWarnParam(c *hbtp.Context, m *struct {
	Name  string `json:"name"`
	Value *comm.ParamWarn
}) {
	if m.Name == "" || m.Value.WarnVal <= 0 {
		c.ResString(hbtp.ResStatusErr, "param err")
		return
	}

	var err error
	switch m.Name {
	case "warn-cpu-avg", "warn-mem-per", "warn-swap-per":
		err = service.SetsParam(m.Name, m.Value)
	default:
		err = errors.New("Not Found Name in [warn-cpu-avg,warn-mem-per,warn-swap-per]")
	}
	if err != nil {
		c.ResString(hbtp.ResStatusErr, "set param err:"+err.Error())
		return
	}
	Mgr.clearTmr()
	c.ResString(hbtp.ResStatusOk, "ok")
}

func (cmds) GetWarns(c *hbtp.Context) {
	var ls []*model.ReportWarn
	err := comm.Db.OrderBy("id DESC").Find(&ls)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, "find err:"+err.Error())
		return
	}
	c.ResJson(hbtp.ResStatusOk, ls)
}

func (cmds) GetInfos(c *hbtp.Context) {
	var ls []*model.ReportInfo
	err := comm.Db.OrderBy("id DESC").Find(&ls)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, "find err:"+err.Error())
		return
	}
	c.ResJson(hbtp.ResStatusOk, ls)
}
