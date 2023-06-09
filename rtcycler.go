package rtc

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// RTCycler defines the parameters
type RTCycler struct {
	// AppDirName application's directory name (required)
	AppDirName string

	// Config application's configuration (required)
	Config Config

	// AppName application's name (default: app)
	AppName string

	// ConfigFilename configuration's filename (default: app.json)
	ConfigFilename string

	// NoParseArgs if true, os.Args will not be parsed
	NoParseArgs bool

	// WithDaemon if true, the --daemon flag won't be ignored if
	// passed as an argument
	WithDaemon bool

	// CacheSubdirs list of cache subdirs to be created
	CacheSubdirs []string
	// ConfigSubdirs list of config subdirs to be created
	ConfigSubdirs []string
	// DataSubdirs list of data subdirs to be created
	DataSubdirs []string
	// RuntimeSubdirs list of run subdirs to be created
	RuntimeSubdirs []string
}

func setEnv(rt *RTCycler) {
	if rt.AppDirName == "" {
		panic("RTCycler.AppDirName is required")
	}

	if rt.Config == nil {
		panic("RTCycler.Config is required")
	}

	appDirName = rt.AppDirName
	conf = rt.Config

	if rt.AppName != "" {
		appName = rt.AppName
	}

	if rt.ConfigFilename != "" {
		configFilename = rt.ConfigFilename
	}

	InstanceSuffix()

	// XDG-related
	dataDir = filepath.Join(xdg.DataHome, appDirName)
	configDir = filepath.Join(xdg.ConfigHome, appDirName)
	cacheDir = filepath.Join(xdg.CacheHome, appDirName)
	runtimeDir = filepath.Join(xdg.RuntimeDir, appDirName)
	daemonDir = filepath.Join(varBasePath, appDirName)

	// log-related
	appInstance = filepath.Base(os.Args[0])
	lockFilename = appInstance + ".lock"
	logFilename = appInstance + ".log"
	logFilePath := filepath.Join(dataDir, logFilename)
	logFile = &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    1, // megabytes
		MaxBackups: 5,
		MaxAge:     30,    //days
		Compress:   false, // disabled by default
	}
	log.SetOutput(logFile)
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.ErrorLevel)
}
