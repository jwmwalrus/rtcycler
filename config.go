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
	SetLockFile(lockfile.Lockfile)
	SetDefaults()
}

// SaveConfig Saves application's instance configuration
func SaveConfig(path string) (err error) {
	if conf == nil {
		return fmt.Errorf("there's no instance config available")
	}

	return saveConfig(conf, path, lockFile)
}

// SaveThisConfig Saves the given configuration
func SaveThisConfig(c Config, path string) (err error) {
	if c == nil {
		return fmt.Errorf("configuration is empty, there's nothing to save")
	}

	return saveConfig(c, path, lockFile)
}

func loadConfig(c Config, path, lockFile string) (err error) {
	_, err = os.Stat(path)
	slog.Info("Loading config", "path", path)

	if errors.Is(err, fs.ErrNotExist) {
		if flagUseConfig != "" {
			Fatal("No user-provided configuration file was found", "flag-config", flagUseConfig)
		}
		slog.Info("Configuration filename was not found. Generating one...", "filename", configFilename)
		c.SetFirstRun(true)
		if err = saveConfig(c, path, lockFile); err != nil {
			return
		}
	}

	var lock lockfile.Lockfile
	lock, err = lockfile.New(lockFile)
	if err != nil {
		return
	}

	if err = lock.TryLock(); err != nil {
		return
	}

	defer func() {
		if err := lock.Unlock(); err != nil {
			fmt.Printf("Cannot unlock %q, reason: %v\n", lock, err)
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

	c.SetLockFile(lock)
	return
}

func saveConfig(c Config, path, lockFile string) (err error) {
	c.SetDefaults()
	slog.Info("Saving config", "path", path)

	var lock lockfile.Lockfile
	lock, err = lockfile.New(lockFile)
	if err != nil {
		return
	}

	if err = lock.TryLock(); err != nil {
		return
	}

	defer func() {
		if err := lock.Unlock(); err != nil {
			fmt.Printf("Cannot unlock %q, reason: %v\n", lock, err)
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
