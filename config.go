package rtc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/nightlyone/lockfile"
)

// Config defines the configuration interface
type Config interface {
	GetFirstRun() bool
	SetFirstRun(bool)
	SetDefaults()
}

// ConfigLocker defines the config-lockfile interface
type ConfigLocker interface {
	SetLockFile(lockfile.Lockfile)
}

// ConfigResolver defines the config-resolve interface
type ConfigResolver interface {
	Resolve()
}

// SaveConfig Saves application's instance configuration
func SaveConfig(path string) (err error) {
	if conf == nil {
		return fmt.Errorf("there's no instance config available")
	}

	return saveConfig(conf, path)
}

// SaveThisConfig Saves the given configuration
func SaveThisConfig(c Config, path string) (err error) {
	if c == nil {
		return fmt.Errorf("configuration is empty, there's nothing to save")
	}

	return saveConfig(c, path)
}

func checkConfigFileLock() (err error) {
	if configFileLock != "" {
		return
	}

	configFileLock, err = lockfile.New(lockFile)
	return
}

func loadConfig(c Config, path string) (err error) {
	_, err = os.Stat(path)
	slog.Info("Loading config", "path", path)

	if errors.Is(err, fs.ErrNotExist) {
		if flagUseConfig != "" {
			Fatal("No user-provided configuration file was found", "flag-config", flagUseConfig)
		}
		slog.Info("Configuration filename was not found. Generating one...", "filename", configFilename)
		c.SetFirstRun(true)
		if err = saveConfig(c, path); err != nil {
			return
		}
	}

	if err != nil {
		return
	}

	if err = checkConfigFileLock(); err != nil {
		return
	}

	if err = configFileLock.TryLock(); err != nil {
		return
	}

	defer func() {
		if err := configFileLock.Unlock(); err != nil {
			fmt.Printf("Cannot unlock %q, reason: %v\n", configFileLock, err)
		}
	}()

	f, err := os.Open(path)
	onerror.Fatal(err)
	defer f.Close()

	bv, _ := io.ReadAll(f)

	err = json.Unmarshal(bv, c)
	if err != nil {
		return
	}

	if cl, ok := c.(ConfigLocker); ok {
		cl.SetLockFile(configFileLock)
	}
	return
}

func saveConfig(c Config, path string) (err error) {
	c.SetDefaults()
	slog.Info("Saving config", "path", path)

	if err = checkConfigFileLock(); err != nil {
		return
	}

	if err = configFileLock.TryLock(); err != nil {
		return
	}

	defer func() {
		if err := configFileLock.Unlock(); err != nil {
			fmt.Printf("Cannot unlock %q, reason: %v\n", configFileLock, err)
		}
	}()

	var file []byte
	file, err = json.MarshalIndent(c, "", " ")
	if err != nil {
		return
	}

	err = os.WriteFile(path, file, 0644)
	return
}
