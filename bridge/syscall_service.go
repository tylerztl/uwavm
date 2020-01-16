// Copyright (c) 2019, Baidu.com, Inc. All Rights Reserved.

package bridge

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/BeDreamCoder/uwavm/common/db"
	"github.com/BeDreamCoder/uwavm/contract/go/pb"
)

// SyscallService is the handler of contract syscalls
type SyscallService struct {
	ctxmgr *StateManager
	db     db.Database
}

// NewSyscallService instances a new SyscallService
func NewSyscallService(ctxmgr *StateManager, db db.Database) *SyscallService {
	return &SyscallService{
		ctxmgr: ctxmgr,
		db:     db,
	}
}

// Transfer implements Syscall interface
func (c *SyscallService) Transfer(ctx context.Context, in *pb.TransferRequest) (*pb.TransferResponse, error) {
	return nil, nil
}

// PutObject implements Syscall interface
func (c *SyscallService) PutObject(ctx context.Context, in *pb.PutRequest) (*pb.PutResponse, error) {
	nctx, ok := c.ctxmgr.GetContractState(in.GetHeader().Ctxid)
	if !ok {
		return nil, fmt.Errorf("bad ctx id:%d", in.Header.Ctxid)
	}
	if in.Value == nil {
		return nil, errors.New("put nil value")
	}
	compk := fmt.Sprintf("%s-%s", nctx.ContractName, string(in.Key))
	err := c.db.Put([]byte(compk), in.Value)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to PutObject for key:[%s],value:[%s]", compk, string(in.Value)))
	}

	return &pb.PutResponse{}, nil
}

// GetObject implements Syscall interface
func (c *SyscallService) GetObject(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	nctx, ok := c.ctxmgr.GetContractState(in.GetHeader().Ctxid)
	if !ok {
		return nil, fmt.Errorf("bad ctx id:%d", in.Header.Ctxid)
	}
	compk := fmt.Sprintf("%s-%s", nctx.ContractName, string(in.Key))
	value, err := c.db.Get([]byte(compk))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Cant GetObject for key: [%s]", compk))
	}
	return &pb.GetResponse{
		Value: value,
	}, nil
}

// DeleteObject implements Syscall interface
func (c *SyscallService) DeleteObject(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	nctx, ok := c.ctxmgr.GetContractState(in.GetHeader().Ctxid)
	if !ok {
		return nil, fmt.Errorf("bad ctx id:%d", in.Header.Ctxid)
	}
	compk := fmt.Sprintf("%s-%s", nctx.ContractName, string(in.Key))
	err := c.db.Delete([]byte(compk))
	return &pb.DeleteResponse{}, err
}

// GetCallArgs implements Syscall interface
func (c *SyscallService) GetCallArgs(ctx context.Context, in *pb.GetCallArgsRequest) (*pb.CallArgs, error) {
	nctx, ok := c.ctxmgr.GetContractState(in.GetHeader().Ctxid)
	if !ok {
		return nil, fmt.Errorf("bad ctx id:%d", in.Header.Ctxid)
	}
	var args []*pb.ArgPair
	for key, value := range nctx.Args {
		args = append(args, &pb.ArgPair{
			Key:   key,
			Value: value,
		})
	}
	sort.Slice(args, func(i, j int) bool {
		return args[i].Key < args[j].Key
	})
	return &pb.CallArgs{
		Method:    nctx.Method,
		Args:      args,
		Caller: nctx.Caller,
	}, nil
}

// SetOutput implements Syscall interface
func (c *SyscallService) SetOutput(ctx context.Context, in *pb.SetOutputRequest) (*pb.SetOutputResponse, error) {
	nctx, ok := c.ctxmgr.GetContractState(in.Header.Ctxid)
	if !ok {
		return nil, fmt.Errorf("bad ctx id:%d", in.Header.Ctxid)
	}
	nctx.Output = in.GetResponse()
	return new(pb.SetOutputResponse), nil
}
