package server

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/yggworldtree/cpu_report/comm"
	"github.com/yggworldtree/cpu_report/model"
	"github.com/yggworldtree/cpu_report/service"
	"github.com/yggworldtree/go-core/bean"
	"github.com/yggworldtree/go-core/utils"
)

type Manager struct {
	Ctx    context.Context
	cncl   context.CancelFunc
	reging bool

	tmr *utils.Timer

	blk    sync.RWMutex
	cpuDev *bean.CliGroupPath

	wrking   bool
	uppartmr *utils.Timer

	warnInterval time.Duration
	warnCpuAvg   []comm.ParamWarn
	warnMemPer   []comm.ParamWarn
	warnSwapPer  []comm.ParamWarn

	warnTmlk sync.RWMutex
	warnTmrs map[string]*utils.Timer
}

func NewManager() *Manager {
	c := &Manager{
		tmr:      utils.NewTimer(time.Hour),
		uppartmr: utils.NewTimer(time.Second * 30),
		warnTmrs: make(map[string]*utils.Timer),
	}
	c.Ctx, c.cncl = context.WithCancel(context.Background())
	return c
}

func (c *Manager) init() error {
	c.warnInterval = time.Second * 30
	bts, err := service.GetParam("warn-interval")
	if err != nil {
		bts = []byte(fmt.Sprintf("%d", c.warnInterval))
		err = service.SetParam("warn-interval", bts)
		if err != nil {
			return err
		}
	} else {
		n, _ := strconv.ParseInt(string(bts), 10, 64)
		if n > 0 {
			hbtp.Debugf("warnInterval:%d", n)
			c.warnInterval = time.Duration(n)
		}
	}
	err = service.GetsParam("warn-cpu-avg", &c.warnCpuAvg)
	if err != nil {
		c.warnCpuAvg = append(c.warnCpuAvg, comm.ParamWarn{
			WarnVal:  80,
			WarnTips: "CpuAvg {{value}}>{{wvalue}}!!",
		})
		c.warnCpuAvg = append(c.warnCpuAvg, comm.ParamWarn{
			WarnVal:  90,
			WarnTips: "CpuAvg {{value}}>{{wvalue}}!!!",
		})
		err = service.SetsParam("warn-cpu-avg", &c.warnCpuAvg)
		if err != nil {
			return err
		}
	}
	err = service.GetsParam("warn-mem-per", &c.warnMemPer)
	if err != nil {
		c.warnMemPer = append(c.warnMemPer, comm.ParamWarn{
			WarnVal:  80,
			WarnTips: "MemPercent {{value}}>{{wvalue}}!!",
		})
		c.warnMemPer = append(c.warnMemPer, comm.ParamWarn{
			WarnVal:  90,
			WarnTips: "MemPercent {{value}}>{{wvalue}}!!!",
		})
		err = service.SetsParam("warn-mem-per", &c.warnMemPer)
		if err != nil {
			return err
		}
	}

	for i, v := range c.warnCpuAvg {
		hbtp.Debugf("CpuAvg[%d]:%.4f", i, v.WarnVal)
	}

	go func() {
		for !hbtp.EndContext(c.Ctx) {
			c.run()
			time.Sleep(time.Second)
		}
	}()

	return nil
}
func (c *Manager) clearTmr() {
	c.warnTmlk.Lock()
	defer c.warnTmlk.Unlock()
	c.warnTmrs = make(map[string]*utils.Timer)
	c.uppartmr.Reset(false)
}
func (c *Manager) run() {
	defer func() {
		if err := recover(); err != nil {
			hbtp.Debugf("Manager run recover:%v", err)
		}
	}()

	if c.reging {
		return
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	lastday := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
	/* if now.Hour() != 1 && now.Hour() != 2 {
		return
	} */
	if !c.tmr.Tick() {
		return
	}

	n, err := comm.Db.Where("start_time>=?", lastday.Format(comm.TimeFmts)).
		Count(&model.ReportInfo{})
	if err != nil || n > 0 {
		hbtp.Debugf("warnInfo lastday had:%v", err)
		return
	}

	n, err = comm.Db.Where("create_time>=? and create_time<?",
		lastday.Format(comm.TimeFmts), today.Format(comm.TimeFmts)).
		Count(&model.ReportWarn{})
	if err != nil {
		hbtp.Debugf("warn count err:%v", err)
		return
	}

	ne := &model.ReportInfo{}
	ne.CreateTime = time.Now()
	ne.UpdateTime = ne.CreateTime
	ne.Type = model.InfoTypeDay
	ne.StartTime = lastday
	ne.EndTime = today
	ne.WarnLen = int(n)
	_, err = comm.Db.InsertOne(ne)
	if err != nil {
		hbtp.Errorf("insert ReportInfo err:%v", err)
		return
	}
}

func (c *Manager) startReg() {
	defer func() {
		if err := recover(); err != nil {
			hbtp.Debugf("Manager StartReg recover:%v", err)
		}
	}()
	if c.reging {
		return
	}
	c.reging = true
	defer func() {
		c.reging = false
	}()
	for !hbtp.EndContext(c.Ctx) {
		err := YwtEgn.SubTopic(comm.MsgPthCpuMem)
		if err == nil {
			break
		}
		hbtp.Debugf("SubTopic %s err:%v", comm.MsgPthCpuMem.String(), err)
		time.Sleep(time.Second)
	}
}
func (c *Manager) startWrk(box *comm.MsgBox) {
	defer func() {
		if err := recover(); err != nil {
			hbtp.Debugf("Manager StartReg recover:%v", err)
		}
	}()
	if c.wrking {
		return
	}
	c.wrking = true
	defer func() {
		c.wrking = false
	}()
	hbtp.Debugf("start update db")

	ne := &model.ReportLog{
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		Name:       box.Name,
		CpuAvg:     box.Cpu.Average,
		ProcLen:    box.Cpu.ProcessLen,
		MemTotal:   int64(box.VirtualMem.Total),
		SwapTotal:  int64(box.SwapMem.Total),
		MemUsed:    int64(box.VirtualMem.Used),
		SwapUsed:   int64(box.SwapMem.Used),
		MemPer:     box.VirtualMem.UsedPercent,
		SwapPer:    box.SwapMem.UsedPercent,
	}
	_, err := comm.Db.InsertOne(ne)
	if err != nil {
		hbtp.Debugf("report log insert err:%v", err)
	}

	if c.uppartmr.Tick() {
		service.GetsParam("warn-cpu-avg", &c.warnCpuAvg)
		service.GetsParam("warn-mem-per", &c.warnMemPer)
		service.GetsParam("warn-swap-per", &c.warnSwapPer)
	}

	c.outCpuWarns(box.Cpu.Average, c.warnCpuAvg)
	c.outCpuWarns(box.VirtualMem.UsedPercent, c.warnMemPer)
	c.outCpuWarns(box.SwapMem.UsedPercent, c.warnSwapPer)
}
func (c *Manager) outCpuWarns(val float64, ws []comm.ParamWarn) {
	gid := utils.NewXid()
	for i, v := range ws {
		k := fmt.Sprintf("%s:%d", "cpu-avg", i)
		c.warnTmlk.RLock()
		tmr := c.warnTmrs[k]
		c.warnTmlk.RUnlock()
		if tmr == nil {
			tmr = utils.NewTimer(c.warnInterval)
			c.warnTmlk.Lock()
			c.warnTmrs[k] = tmr
			c.warnTmlk.Unlock()
		}
		if v.WarnVal > 0 && val >= v.WarnVal && tmr.Tick() {
			c.outWarn(gid, val, v.WarnVal, i+1, v.WarnTips)
		}
	}
}
func (c *Manager) outWarn(gid string, v, wv float64, lev int, warns string) {
	warns = strings.ReplaceAll(warns, "{{value}}", fmt.Sprintf("%.4f", v))
	warns = strings.ReplaceAll(warns, "{{wvalue}}", fmt.Sprintf("%.4f", wv))
	nw := &model.ReportWarn{
		Gid:        gid,
		CreateTime: time.Now(),
		Type:       model.WarnTypeCpuAvgOut,
		Val:        v,
		Wval:       wv,
		Lev:        lev,
		Warns:      warns,
	}
	_, err := comm.Db.InsertOne(nw)
	if err != nil {
		hbtp.Debugf("report warn insert err:%v", err)
	}
}
