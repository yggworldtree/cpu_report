package model

import (
	"time"
)

type ReportInfo struct {
	Id         int64     `xorm:"pk autoincr BIGINT(20)" json:"id"`
	CreateTime time.Time `xorm:"DATETIME" json:"createTime"`
	UpdateTime time.Time `xorm:"DATETIME" json:"updateTime"`
	Type       string    `xorm:"VARCHAR(50)" json:"type"`
	StartTime  time.Time `xorm:"DATETIME" json:"startTime"`
	EndTime    time.Time `xorm:"DATETIME" json:"endTime"`
	Infos      string    `xorm:"LONGTEXT" json:"infos"`
}
