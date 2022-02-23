package rpc

import (
	"bytes"
	"context"
	"deck/service/db/mapper"
	"deck/service/db/orm"
	dbProxy "deck/service/db/proto"
	"encoding/json"
)

type DB struct{}

func (db *DB) ExecuteAction(ctx context.Context, req *dbProxy.ReqExec, res *dbProxy.RespExec) error {
	resList := make([]orm.ExecResult, len(req.Action))
	for idx, singleAciton := range req.Action {
		params := []interface{}{}
		dec := json.NewDecoder(bytes.NewReader(singleAciton.Params))
		dec.UseNumber()
		if err := dec.Decode(&params); err != nil {
			resList[idx] = orm.ExecResult{
				Suc: false,
				Msg: "请求参数有误",
			}
			continue
		}
		for k, v := range params {
			if _, ok := v.(json.Number); ok {
				params[k], _ = v.(json.Number).Int64()
			}
		}
		execRes, err := mapper.FuncCall(singleAciton.Name, params...)
		if err != nil {
			resList[idx] = orm.ExecResult{
				Suc: false,
				Msg: "函数点用有误",
			}
			continue
		}
		resList[idx] = execRes[0].Interface().(orm.ExecResult)
	}
	res.Data, _ = json.Marshal(resList)
	return nil
}
