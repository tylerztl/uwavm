package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/BeDreamCoder/uwavm/bridge"
	"github.com/BeDreamCoder/uwavm/common/db/leveldb"
	"github.com/BeDreamCoder/uwavm/vm"
)

var (
	centry  = flag.String("entry", "run", "entry function")
	environ = flag.String("e", "go", "environ, c or go")
)

func makeDeployArgs(modulePath string) map[string][]byte {
	codebuf, err := ioutil.ReadFile(modulePath)
	if err != nil {
		panic(err)
	}

	args := map[string][]byte{
		"initSupply": []byte("1000000"),
	}
	argsbuf, _ := json.Marshal(args)
	return map[string][]byte{
		"contract_name": []byte("erc20"),
		"contract_code": codebuf,
		"language":      []byte("go"),
		"args":          argsbuf,
		"caller":        []byte("alice"),
	}
}

func makeInvokeArgs() map[string][]byte {
	args := map[string][]byte{
		"action":  []byte("balanceOf"),
		"address": []byte("alice"),
	}
	argsbuf, _ := json.Marshal(args)
	return map[string][]byte{
		"contract_name": []byte("erc20"),
		"language":      []byte("go"),
		"args":          argsbuf,
		"caller":        []byte("alice"),
	}
}

func run(modulePath string) error {
	_, err := filepath.Abs(modulePath)
	if err != nil {
		return err
	}

	db := leveldb.NewProvider().GetDBHandle("uwavm")
	bridge := bridge.NewBridge(db)
	vm := vm.NewVMManager(db, bridge)
	bridge.RegisterExecutor("wasm", vm)
	resp, err := vm.DeployContract(makeDeployArgs(modulePath))
	//resp, err := vm.InvokeContract("query", makeInvokeArgs())
	fmt.Println("Status:", resp.GetStatus())
	fmt.Println("Message:", resp.GetMessage())
	fmt.Println("Bdoy:", string(resp.GetBody()))
	return err
}

func main() {
	flag.Parse()

	filename := flag.Arg(0)
	ext := filepath.Ext(filename)
	var target string
	var err error
	switch ext {
	case ".wasm":
		target = flag.Arg(0)
	default:
		log.Fatalf("bad file ext:%s", ext)
	}

	err = run(target)
	if err != nil {
		log.Fatal(err)
	}
}
