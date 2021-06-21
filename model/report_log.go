package model

import (
	"time"
)

type ReportLog struct {
	Id         int64     `xorm:"pk autoincr comment('primary key') BIGINT(20)" json:"id"`
	CreateTime time.Time `xorm:"comment('create time') DATETIME" json:"createTime"`
	UpdateTime time.Time `xorm:"comment('update time') DATETIME" json:"updateTime"`
	Name       string    `xorm:"comment('dev name') VARCHAR(255)" json:"name"`
	CpuAvg     float64   `xorm:"DOUBLE" json:"cpuAvg"`
	ProcLen    int       `xorm:"INT(11)" json:"procLen"`
	MemTotal   int64     `xorm:"BIGINT(20)" json:"memTotal"`
	SwapTotal  int64     `xorm:"BIGINT(20)" json:"swapTotal"`
	MemUsed    int64     `xorm:"BIGINT(20)" json:"memUsed"`
	SwapUsed   int64     `xorm:"BIGINT(20)" json:"swapUsed"`
	MemPer     float64   `xorm:"DOUBLE" json:"memPer"`
	SwapPer    float64   `xorm:"DOUBLE" json:"swapPer"`
}
