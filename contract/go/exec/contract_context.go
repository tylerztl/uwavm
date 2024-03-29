package exec

import (
	"log"

	"github.com/BeDreamCoder/uwavm/contract/go/code"
	"github.com/BeDreamCoder/uwavm/contract/go/pb"
)

const (
	methodPut          = "PutObject"
	methodGet          = "GetObject"
	methodDelete       = "DeleteObject"
	methodOutput       = "SetOutput"
	methodGetCallArgs  = "GetCallArgs"
	methodContractCall = "ContractCall"
)

type contractContext struct {
	callArgs       pb.CallArgs
	contractArgs   map[string][]byte
	bridgeCallFunc BridgeCallFunc
	header         pb.SyscallHeader
}

func newContractContext(ctxid int64, bridgeCallFunc BridgeCallFunc) *contractContext {
	return &contractContext{
		contractArgs:   make(map[string][]byte),
		bridgeCallFunc: bridgeCallFunc,
		header: pb.SyscallHeader{
			Ctxid: ctxid,
		},
	}
}

func (c *contractContext) Init() error {
	var request pb.GetCallArgsRequest
	request.Header = &c.header
	err := c.bridgeCallFunc(methodGetCallArgs, &request, &c.callArgs)
	if err != nil {
		return err
	}
	for _, pair := range c.callArgs.GetArgs() {
		c.contractArgs[pair.GetKey()] = pair.GetValue()
	}
	return nil
}

func (c *contractContext) Method() string {
	return c.callArgs.GetMethod()
}

func (c *contractContext) Args() map[string][]byte {
	return c.contractArgs
}

func (c *contractContext) Caller() string {
	return c.callArgs.Caller
}

func (c *contractContext) PutObject(key, value []byte) error {
	req := &pb.PutRequest{
		Header: &c.header,
		Key:    key,
		Value:  value,
	}
	rep := new(pb.PutResponse)
	return c.bridgeCallFunc(methodPut, req, rep)
}

func (c *contractContext) GetObject(key []byte) ([]byte, error) {
	req := &pb.GetRequest{
		Header: &c.header,
		Key:    key,
	}
	rep := new(pb.GetResponse)
	err := c.bridgeCallFunc(methodGet, req, rep)
	if err != nil {
		return nil, err
	}
	return rep.Value, nil
}

func (c *contractContext) DeleteObject(key []byte) error {
	req := &pb.DeleteRequest{
		Header: &c.header,
		Key:    key,
	}
	rep := new(pb.DeleteResponse)
	return c.bridgeCallFunc(methodDelete, req, rep)
}

func (c *contractContext) Call(module, contract, method string, args map[string][]byte) (*code.Response, error) {
	var argPairs []*pb.ArgPair
	// 在合约里面单次合约调用的map迭代随机因子是确定的，因此这里不需要排序
	for key, value := range args {
		argPairs = append(argPairs, &pb.ArgPair{
			Key:   key,
			Value: value,
		})
	}
	req := &pb.ContractCallRequest{
		Header:   &c.header,
		Module:   module,
		Contract: contract,
		Method:   method,
		Args:     argPairs,
	}
	rep := new(pb.ContractCallResponse)
	err := c.bridgeCallFunc(methodContractCall, req, rep)
	if err != nil {
		return nil, err
	}
	return &code.Response{
		Status:  int(rep.Response.Status),
		Message: rep.Response.Message,
		Body:    rep.Response.Body,
	}, nil
}

func (c *contractContext) SetOutput(response *code.Response) error {
	req := &pb.SetOutputRequest{
		Header: &c.header,
		Response: &pb.Response{
			Status:  int32(response.Status),
			Message: response.Message,
			Body:    response.Body,
		},
	}
	rep := new(pb.SetOutputResponse)
	err := c.bridgeCallFunc(methodOutput, req, rep)
	if err != nil {
		log.Printf("Setoutput error:%s", err)
	}
	return err
}
