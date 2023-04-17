package rtc

import (
	"path/filepath"

	"github.com/jwmwalrus/bnp/env"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
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

	log.WithField("flagDaemonMode", flagDaemonMode).
		Infof("Daemon mode allowed: %v", rt.WithDaemon)
	if rt.WithDaemon && flagDaemonMode {
		log.WithField("daemonDir", daemonDir).
			Info("Applying daemon mode")
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
	onerror.Panic(err)

	configFile = filepath.Join(configDir, configFilename)
	if flagUseConfig != "" {
		configFile = flagUseConfig
		log.Info("Using provided config file instead")
	}
	log.Infof("Using config file: %s", configFile)

	lockFile = filepath.Join(runtimeDir, lockFilename)
	log.Infof("Using lock file: %s", lockFile)

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
		onerror.Panic(err)
	}

	err = loadConfig(conf, configFile, lockFile)
	onerror.Panic(err)

	if flagUseConfig != "" {
		configFile = filepath.Join(configDir, configFilename)
		log.WithField("configFile", configFile).
			Info("Restored config file to its default")
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
	log.Info("Unloading application")

	if conf == nil {
		return
	}

	if conf.GetFirstRun() {
		conf.SetFirstRun(false)

		err := saveConfig(conf, configFile, lockFile)
		onerror.Log(err)
	}

	for i := len(unloadRegistry) - 1; i >= 0; i-- {
		log.Infof("Calling unloader: %v", unloadRegistry[i].Description)
		onerror.Log(unloadRegistry[i].Callback())
	}
}
