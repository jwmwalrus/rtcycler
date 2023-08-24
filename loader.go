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
func Load(rt RTCycler) (args []string) {
	setEnv(&rt)

	if !rt.NoParseArgs {
		args = parseArgs()
	}

	slog.With(
		"flag-daemon-mode", flagDaemonMode,
		"allow-daemon-mode", rt.WithDaemon,
	).Info("Daemon mode status")
	if rt.WithDaemon && flagDaemonMode {
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
	for i := range rt.CacheSubdirs {
		list = append(list, filepath.Join(cacheDir, rt.CacheSubdirs[i]))
	}
	for i := range rt.ConfigSubdirs {
		list = append(list, filepath.Join(configDir, rt.ConfigSubdirs[i]))
	}
	for i := range rt.DataSubdirs {
		list = append(list, filepath.Join(dataDir, rt.DataSubdirs[i]))
	}
	for i := range rt.RuntimeSubdirs {
		list = append(list, filepath.Join(runtimeDir, rt.RuntimeSubdirs[i]))
	}
	if len(list) > 0 {
		err = env.CreateTheseDirs(list)
		onerror.Fatal(err)
	}

	err = loadConfig(conf, configFile, lockFile)
	onerror.Fatal(err)

	if flagUseConfig != "" {
		configFile = filepath.Join(configDir, configFilename)
		slog.Info("Restored configuration file to its default", "config-file", configFile)
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

	if conf.GetFirstRun() {
		conf.SetFirstRun(false)

		err := saveConfig(conf, configFile, lockFile)
		onerror.Log(err)
	}

	for i := len(unloadRegistry) - 1; i >= 0; i-- {
		slog.Info("Calling unloader: %v", "unloader-descr", unloadRegistry[i].Description)
		onerror.Log(unloadRegistry[i].Callback())
	}
}
