package service

import (
	"encoding/json"
	"errors"
	"time"

	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/yggworldtree/cpu_report/comm"
	"github.com/yggworldtree/cpu_report/model"
)

func FindParam(key string) (*model.ReportParam, bool) {
	e := &model.ReportParam{}
	ok, err := comm.Db.Where("name=?", key).Get(e)
	if err != nil {
		hbtp.Debugf("FindParam(%s) err:%v", key, err)
	}
	return e, ok
}
func SetParam(key string, data []byte, tit ...string) error {
	var err error
	e, ok := FindParam(key)
	if len(tit) > 0 {
		e.Title = tit[0]
	}
	e.Data = string(data)
	e.UpdateTime = time.Now()
	if ok && e.Id > 0 {
		_, err = comm.Db.Cols("title", "data").Where("id=?", e.Id).Update(e)
	} else {
		e.Name = key
		e.CrateTime = time.Now()
		_, err = comm.Db.Insert(e)
	}
	return err
}
func SetsParam(key string, data interface{}, tit ...string) error {
	if data == nil {
		return errors.New("data is nil")
	}
	bts, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return SetParam(key, bts, tit...)
}

func GetParam(key string) ([]byte, error) {
	e, ok := FindParam(key)
	if ok {
		return []byte(e.Data), nil
	}
	return nil, errors.New("not found param")
}
func GetsParam(key string, data interface{}) error {
	if data == nil {
		return errors.New("data is nil")
	}
	bts, err := GetParam(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(bts, data)
}
