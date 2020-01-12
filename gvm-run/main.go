package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/BeDreamCoder/uwavm/debug"
	"github.com/BeDreamCoder/uwavm/exec"
	gowasm "github.com/BeDreamCoder/uwavm/runtime/go"
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

func run(modulePath string, args []string) error {
	_, err := filepath.Abs(modulePath)
	if err != nil {
		return err
	}
	resolver := exec.NewMultiResolver(resolver, gowasm.NewResolver())
	var code exec.Code
	codebuf, err := ioutil.ReadFile(modulePath)
	if err != nil {
		return err
	}
	code, err = exec.NewInterpCode(codebuf, resolver)
	if err != nil {
		return err
	}

	defer code.Release()
	ctx, err := code.NewContext(exec.DefaultContextConfig())
	if err != nil {
		return err
	}

	defer ctx.Release()
	debug.SetWriter(ctx, os.Stderr)
	var entry string
	switch *environ {
	case "go":
		entry = "run"
		gowasm.RegisterRuntime(ctx)
	}

	var argc, argv int
	if ctx.Memory() != nil {
		argc, argv = prepareArgs(ctx.Memory(), args, nil)
	}
	ret, err := ctx.Exec(entry, []int64{int64(argc), int64(argv)})
	fmt.Println("gas: ", ctx.GasUsed())
	fmt.Println("ret: ", ret)
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

	err = run(target, flag.Args()[0:])
	if err != nil {
		log.Fatal(err)
	}
}
