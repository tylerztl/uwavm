package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/BeDreamCoder/uwavm/bridge"
	"github.com/BeDreamCoder/uwavm/common/db/leveldb"
	"github.com/BeDreamCoder/uwavm/wasm"
)

var (
	centry  = flag.String("entry", "run", "entry function")
	environ = flag.String("e", "go", "environ, c or go")
)

func prepareArgs(mem []byte, args []string, envs []string) (int, int) {
	argc := len(args)
	offset := 4 << 10
	strdup := func(s string) int {
		copy(mem[offset:], s+"\x00")
		ptr := offset
		offset += len(s) + (8 - len(s)%8)
		return ptr
	}
	var argvAddr []int
	for _, arg := range args {
		argvAddr = append(argvAddr, strdup(arg))
	}

	argvAddr = append(argvAddr, len(envs))
	for _, env := range envs {
		argvAddr = append(argvAddr, strdup(env))
	}

	argv := offset
	buf := bytes.NewBuffer(mem[offset:offset])
	for _, addr := range argvAddr {
		if *environ == "go" {
			binary.Write(buf, binary.LittleEndian, uint64(addr))
		} else {
			binary.Write(buf, binary.LittleEndian, uint32(addr))
		}
	}
	return argc, argv
}

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
		"init_args":     argsbuf,
	}
}

func run(modulePath string) error {
	_, err := filepath.Abs(modulePath)
	if err != nil {
		return err
	}

	db := leveldb.NewProvider().GetDBHandle("uwasm")
	bridge := bridge.NewBridge(db)
	vm := wasm.NewVMManager(db, bridge)
	bridge.RegisterExecutor("uwasm", vm)
	resp, err := vm.DeployContract(makeDeployArgs(modulePath))

	fmt.Println("Status:", resp.GetStatus())
	fmt.Println("Message:", resp.GetMessage())
	fmt.Println("Bdoy:", resp.GetBody())

	//cctx.Method = "invoke"
	//cctx.Args = map[string][]byte{"action": []byte("transfer"),
	//	"to":     []byte("bob"),
	//	"amount": []byte("1")}
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
