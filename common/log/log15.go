package log

import (
	"go/build"
	"os"
	"path"
	"path/filepath"
	"sync"

	log "github.com/inconshreveable/log15"
)

var (
	filePath, errorFile string
	logOnce             sync.Once
	singleton           log.Logger
)

func init() {
	filePath = path.Join(goPath(), "src/github.com/BeDreamCoder/uwavm")
	errorFile = path.Join(filePath, "error.json")
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		panic(err)
	}
}

func goPath() string {
	gpDefault := build.Default.GOPATH
	gps := filepath.SplitList(gpDefault)

	return gps[0]
}

func New(ctx ...interface{}) log.Logger {
	uLog := log.New(ctx)
	uLog.SetHandler(log.SyncHandler(log.MultiHandler(
		log.StreamHandler(os.Stderr, log.LogfmtFormat()),
		log.LvlFilterHandler(log.LvlError, log.Must.FileHandler(
			errorFile, log.JsonFormat())))))

	return uLog
}

func GetLogger() log.Logger {
	logOnce.Do(func() {
		singleton = New("uwasm")
	})
	return singleton
}
