package handler

import (
	"context"
	"deck/common"
	"deck/service/account/proto"
	dbcli "deck/service/db/client"
	"encoding/json"
)

func (u *User) UserFiles(ctx context.Context, req *proto.ReqUserFile, res *proto.RespUserFile) error {
	dbResp, err := dbcli.QueryUserFileMetas(req.Username, int(req.Limit))
	if err != nil || !dbResp.Suc {
		res.Code = common.StatusServerError
		return err
	}
	userFiles := dbcli.ToTableUserFiles(dbResp.Data)
	data, err := json.Marshal(userFiles)
	if err != nil {
		res.Code = common.StatusServerError
		return nil
	}
	res.FileData = data
	return nil
}
func (u *User) UserFileRename(ctx context.Context, req *proto.ReqUserFileRename, res *proto.RespUserFileRename) error {
	dbResp, err := dbcli.RenameFileName(req.Username, req.Filehash, req.NewFileName)
	if err != nil || dbResp.Suc {
		res.Code = common.StatusServerError
		return nil
	}
	userFiles := dbcli.ToTableUserFiles(dbResp.Data)
	data, err := json.Marshal(userFiles)
	if err != nil {
		res.Code = common.StatusServerError
		return nil
	}
	res.FileData = data
	return nil
}
