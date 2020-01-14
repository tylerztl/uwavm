package log

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
)

var logger *logs.BeeLogger

func init() {
	logger = logs.NewLogger()
	config := make(map[string]interface{})
	logPath := "logs"
	if fi, err := os.Stat(logPath); err != nil {
		if err := os.MkdirAll(logPath, 0755); err != nil {
			panic("Invalid log path")
		}
	} else if !fi.IsDir() {
		panic(fmt.Sprintf("%s must be a directory", logPath))
	}
	logFile := filepath.Join(logPath, "logagent.log")
	if _, err := os.Stat(logFile); err != nil {
		if err = ioutil.WriteFile(logFile, nil, 0644); err != nil {
			panic(err)
		}
	}
	config["filename"] = logFile
	config["level"] = logs.LevelDebug
	configStr, err := json.Marshal(config)
	if err != nil {
		panic(fmt.Errorf("logger marshal err[%s]", err))
	}
	err = logger.SetLogger(logs.AdapterConsole, string(configStr))
	err = logger.SetLogger(logs.AdapterFile, string(configStr))
	if err != nil {
		panic(fmt.Errorf("logger SetLogger err[%s]", err))
	}
	logger.EnableFuncCallDepth(true)
	logger.SetLogFuncCallDepth(4)
}

func GetLogger() *logs.BeeLogger {
	return logger
}
