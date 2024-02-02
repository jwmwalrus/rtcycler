package rtc

import (
	"log/slog"
	"path/filepath"

	"github.com/jwmwalrus/bnp/env"
	"github.com/jwmwalrus/bnp/onerror"
)

// Unloader defines a method to be called when unloading the application
type Unloader struct {
	Description string
	Callback    UnloaderCallback
}

// UnloaderCallback defines the signature of the method to be called when
// unloading the application
type UnloaderCallback func() error

var (
	unloadRegistry = []*Unloader{}
)

// Load Loads application's configuration
func Load(c Config, appDirName string, options ...func(*RTCycler)) (args []string) {
	rt := &RTCycler{
		appDirName: appDirName,
		config:     c,
	}
	for _, o := range options {
		o(rt)
	}
	setEnv(rt)

	if !rt.noParseArgs {
		args = parseArgs()
	}

	slog.With(
		"flag-daemon-mode", flagDaemonMode,
		"allow-daemon-mode", rt.daemon,
	).Info("Daemon mode status")
	if rt.daemon && flagDaemonMode {
		slog.Info("Applying daemon mode", "daemon-dir", daemonDir)
		cacheDir = daemonDir
		configDir = daemonDir
		dataDir = daemonDir
		runtimeDir = daemonDir
	}

	err := env.SetDirs(
		cacheDir,
		configDir,
		dataDir,
		runtimeDir,
	)
	onerror.Fatal(err)

	configFile = filepath.Join(configDir, configFilename)
	if flagUseConfig != "" {
		configFile = flagUseConfig
		slog.Info("Using provided config file instead")
	}
	slog.Info("Using configuration file", "config-file", configFile)

	lockFile = filepath.Join(runtimeDir, lockFilename)
	slog.Info("Using lock file", "lock-file", lockFile)

	var list []string
	for i := range rt.cacheSubdirs {
		list = append(list, filepath.Join(cacheDir, rt.cacheSubdirs[i]))
	}
	for i := range rt.configSubdirs {
		list = append(list, filepath.Join(configDir, rt.configSubdirs[i]))
	}
	for i := range rt.dataSubdirs {
		list = append(list, filepath.Join(dataDir, rt.dataSubdirs[i]))
	}
	for i := range rt.runtimeSubdirs {
		list = append(list, filepath.Join(runtimeDir, rt.runtimeSubdirs[i]))
	}
	if len(list) > 0 {
		err = env.CreateTheseDirs(list)
		onerror.Fatal(err)
	}

	err = loadConfig(conf, configFile)
	onerror.Fatal(err)

	if flagUseConfig != "" {
		configFile = filepath.Join(configDir, configFilename)
		slog.Info("Restored configuration file to its default", "config-file", configFile)
	}

	if cr, ok := rt.config.(ConfigResolver); ok {
		cr.Resolve()
	}
	return
}

// RegisterUnloader registers an Unloader, to be invoked before stopping the app
func RegisterUnloader(u *Unloader) {
	unloadRegistry = append(unloadRegistry, u)
}

// Unload Cleans up server before exit.
// Registered unloaders are invoked in reverse order of registration.
func Unload() {
	slog.Info("Unloading application")

	if conf == nil {
		return
	}

	for i := len(unloadRegistry) - 1; i >= 0; i-- {
		if unloadRegistry[i].Callback == nil {
			continue
		}

		slog.Info("Calling unloader: %v", "unloader-descr", unloadRegistry[i].Description)
		onerror.Log(unloadRegistry[i].Callback())
	}

	if conf.GetFirstRun() {
		conf.SetFirstRun(false)

		err := saveConfig(conf, configFile)
		onerror.Log(err)
	}
}
