package rtc

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"gopkg.in/natefinch/lumberjack.v2"
)

// RTCycler defines the parameters
type RTCycler struct {
	// AppDirName application's directory name (required)
	appDirName string

	// Config application's configuration (required)
	config Config

	// AppName application's name (default: app)
	appName string

	// ConfigFilename configuration's filename (default: app.json)
	configFilename string

	// NoParseArgs if true, os.Args will not be parsed
	noParseArgs bool

	// WithDotHome if true, use a single `dot` home dir instead of XDG
	dotHome bool

	// WithDaemon if true, the --daemon flag won't be ignored if
	// passed as an argument. Overrides WithDotHome
	daemon bool

	// CacheSubdirs list of cache subdirs to be created
	cacheSubdirs []string
	// ConfigSubdirs list of config subdirs to be created
	configSubdirs []string
	// DataSubdirs list of data subdirs to be created
	dataSubdirs []string
	// RuntimeSubdirs list of run subdirs to be created
	runtimeSubdirs []string
}

func WithAppName(s string) func(*RTCycler) {
	return func(rt *RTCycler) {
		rt.appName = s
	}
}

func WithConfigFileName(s string) func(*RTCycler) {
	return func(rt *RTCycler) {
		rt.configFilename = s
	}
}

func WithNoParseArgs() func(*RTCycler) {
	return func(rt *RTCycler) {
		rt.noParseArgs = true
	}
}

func WithDotHome() func(*RTCycler) {
	return func(rt *RTCycler) {
		rt.dotHome = true
	}
}

func WithDaemon() func(*RTCycler) {
	return func(rt *RTCycler) {
		rt.daemon = true
	}
}

func WithCacheSubdirs(l []string) func(*RTCycler) {
	return func(rt *RTCycler) {
		rt.cacheSubdirs = l
	}
}

func WithConfigSubdirs(l []string) func(*RTCycler) {
	return func(rt *RTCycler) {
		rt.configSubdirs = l
	}
}

func WithDataSubdirs(l []string) func(*RTCycler) {
	return func(rt *RTCycler) {
		rt.dataSubdirs = l
	}
}

func WithRuntimeSubdirs(l []string) func(*RTCycler) {
	return func(rt *RTCycler) {
		rt.runtimeSubdirs = l
	}
}

func setEnv(rt *RTCycler) {
	if rt.appDirName == "" {
		panic("RTCycler.AppDirName is required")
	}

	if rt.config == nil {
		panic("RTCycler.Config is required")
	}

	appDirName = rt.appDirName
	conf = rt.config

	if rt.appName != "" {
		appName = rt.appName
	}

	if rt.configFilename != "" {
		configFilename = rt.configFilename
	}

	InstanceSuffix()

	// XDG-related
	if rt.dotHome {
		dotHome := "." + appDirName
		dataDir = filepath.Join(xdg.Home, dotHome)
		configDir = filepath.Join(xdg.Home, dotHome)
		cacheDir = filepath.Join(xdg.Home, dotHome)
	} else {
		dataDir = filepath.Join(xdg.DataHome, appDirName)
		configDir = filepath.Join(xdg.ConfigHome, appDirName)
		cacheDir = filepath.Join(xdg.CacheHome, appDirName)
	}

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

	logLevel.Set(slog.LevelError)
	logger = slog.New(slog.NewTextHandler(logFile, logOptions))
	slog.SetDefault(logger)
}
