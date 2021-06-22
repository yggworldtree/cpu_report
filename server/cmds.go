package server

import (
	"errors"
	"fmt"
	"strconv"
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

func (cmds) GetWarnLen(c *hbtp.Context) {
	days := c.Args().Get("day")
	typs := c.Args().Get("type")
	lev := c.Args().Get("lev")
	day, _ := strconv.Atoi(days)
	ses := comm.Db.NewSession()
	defer ses.Close()
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if day == -1 {
		ses.And("create_time>=?", today.Format(comm.TimeFmts))
	} else if day > 0 {
		day1 := today.AddDate(0, 0, -int(day))
		day2 := today.AddDate(0, 0, -int(day-1))
		ses.And("create_time>=? and create_time<?", day1.Format(comm.TimeFmts), day2.Format(comm.TimeFmts))
	}
	if typs != "" {
		ses.And("`type`=?", typs)
	}
	if lev != "" {
		ses.And("`lev`=?", lev)
	}
	ln, err := ses.Count(&model.ReportWarn{})
	if err != nil {
		c.ResString(hbtp.ResStatusErr, "find err:"+err.Error())
		return
	}
	hbtp.Debugf("GetWarnLen day(%d):%d", day, ln)
	c.ResString(hbtp.ResStatusOk, fmt.Sprintf("%d", ln))
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

func (cmds) GetInfoLen(c *hbtp.Context) {
	ln, err := comm.Db.Count(&model.ReportInfo{})
	if err != nil {
		c.ResString(hbtp.ResStatusErr, "find err:"+err.Error())
		return
	}
	c.ResString(hbtp.ResStatusOk, fmt.Sprintf("%d", ln))
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
