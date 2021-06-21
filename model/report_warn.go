package model

import (
	"time"
)

type ReportWarn struct {
	Id         int64     `xorm:"pk autoincr BIGINT(20)" json:"id"`
	Gid        string    `xorm:"VARCHAR(64)" json:"gid"`
	CreateTime time.Time `xorm:"DATETIME" json:"createTime"`
	Type       string    `xorm:"VARCHAR(50)" json:"type"`
	Val        float64   `xorm:"DOUBLE" json:"val"`
	Wval       float64   `xorm:"DOUBLE" json:"wval"`
	Lev        int       `xorm:"INT(11)" json:"lev"`
	Warns      string    `xorm:"VARCHAR(255)" json:"warns"`
}
