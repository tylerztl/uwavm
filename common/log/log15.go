package log

import (
	"os"
	"path"
	"sync"

	"github.com/BeDreamCoder/uwavm/common/util"
	log "github.com/inconshreveable/log15"
)

var (
	filePath, errorFile string
	logOnce             sync.Once
	singleton           log.Logger
)

func init() {
	filePath = path.Join(util.GoPath(), "src/github.com/BeDreamCoder/uwavm/output/log")
	errorFile = path.Join(filePath, "error.json")
	util.CreateDirIfMissing(filePath)
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
