// Copyright (c) 2019, Baidu.com, Inc. All Rights Reserved.

package bridge

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/BeDreamCoder/uwavm/contract/go/pb"
)

// SyscallService is the handler of contract syscalls
type SyscallService struct {
	ctxmgr *ContextManager
}

// NewSyscallService instances a new SyscallService
func NewSyscallService(ctxmgr *ContextManager) *SyscallService {
	return &SyscallService{
		ctxmgr: ctxmgr,
	}
}

// Transfer implements Syscall interface
func (c *SyscallService) Transfer(ctx context.Context, in *pb.TransferRequest) (*pb.TransferResponse, error) {
	return nil, nil
}

// PutObject implements Syscall interface
func (c *SyscallService) PutObject(ctx context.Context, in *pb.PutRequest) (*pb.PutResponse, error) {
	nctx, ok := c.ctxmgr.Context(in.GetHeader().Ctxid)
	if !ok {
		return nil, fmt.Errorf("bad ctx id:%d", in.Header.Ctxid)
	}
	if in.Value == nil {
		return nil, errors.New("put nil value")
	}
	compk := fmt.Sprintf("%s-%s", nctx.ContractName, string(in.Key))
	ok = nctx.Cache.Add(compk, in.Value)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Failed to PutObject for key:[%s],value:[%s]", compk, string(in.Value)))
	}

	return &pb.PutResponse{}, nil
}

// GetObject implements Syscall interface
func (c *SyscallService) GetObject(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	nctx, ok := c.ctxmgr.Context(in.GetHeader().Ctxid)
	if !ok {
		return nil, fmt.Errorf("bad ctx id:%d", in.Header.Ctxid)
	}
	compk := fmt.Sprintf("%s-%s", nctx.ContractName, string(in.Key))
	value, ok := nctx.Cache.Get(compk)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Cant GetObject for key: [%s]", compk))
	}
	return &pb.GetResponse{
		Value: value.([]byte),
	}, nil
}

// DeleteObject implements Syscall interface
func (c *SyscallService) DeleteObject(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	nctx, ok := c.ctxmgr.Context(in.GetHeader().Ctxid)
	if !ok {
		return nil, fmt.Errorf("bad ctx id:%d", in.Header.Ctxid)
	}
	compk := fmt.Sprintf("%s-%s", nctx.ContractName, string(in.Key))
	nctx.Cache.Del(compk)
	return &pb.DeleteResponse{}, nil
}

// GetCallArgs implements Syscall interface
func (c *SyscallService) GetCallArgs(ctx context.Context, in *pb.GetCallArgsRequest) (*pb.CallArgs, error) {
	nctx, ok := c.ctxmgr.Context(in.GetHeader().Ctxid)
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
		Method:         nctx.Method,
		Args:           args,
		Initiator:      nctx.Initiator,
		AuthRequire:    nctx.AuthRequire,
		TransferAmount: nctx.TransferAmount,
	}, nil
}

// SetOutput implements Syscall interface
func (c *SyscallService) SetOutput(ctx context.Context, in *pb.SetOutputRequest) (*pb.SetOutputResponse, error) {
	nctx, ok := c.ctxmgr.Context(in.Header.Ctxid)
	if !ok {
		return nil, fmt.Errorf("bad ctx id:%d", in.Header.Ctxid)
	}
	nctx.Output = in.GetResponse()
	return new(pb.SetOutputResponse), nil
}
