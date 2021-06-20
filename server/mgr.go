package server

import (
	"context"
	"sync"
	"time"

	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/yggworldtree/cpu_report/comm"
	"github.com/yggworldtree/cpu_report/model"
	"github.com/yggworldtree/go-core/bean"
)

type Manager struct {
	Ctx    context.Context
	cncl   context.CancelFunc
	reging bool

	blk    sync.RWMutex
	cpuDev *bean.CliGroupPath

	wrking bool
}

func NewManager() *Manager {
	c := &Manager{}
	c.Ctx, c.cncl = context.WithCancel(context.Background())
	return c
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
		err := YwtEgn.SubTopic([]*bean.TopicInfo{
			{Path: comm.MsgPthCpuMem.String(), Safed: false},
		})
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
		CpuAvg:     int(box.Cpu.Average),
		ProcLen:    box.Cpu.ProcessLen,
		MemTotal:   int64(box.VirtualMem.Total),
		SwapTotal:  int64(box.SwapMem.Total),
		MemUsed:    int64(box.VirtualMem.Used),
		SwapUsed:   int64(box.SwapMem.Used),
		MemPer:     int(box.VirtualMem.UsedPercent * 4),
		SwapPer:    int(box.SwapMem.UsedPercent * 4),
	}
	_, err := comm.Db.InsertOne(ne)
	if err != nil {
		hbtp.Debugf("report log insert err:%v", err)
	}
}
