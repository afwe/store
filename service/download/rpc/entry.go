package rpc

import (
	cfg "deck/service/download/config"
	dlProto "deck/service/download/proto"
)
import "context"

type Download struct{}

func (u *Download) DownloadEntry(
	ctx context.Context,
	req *dlProto.ReqEntry,
	resp *dlProto.RespEntry,
) error {
	resp.Entry = cfg.DownloadEntry
	return nil
}
