package rtc

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/jwmwalrus/onerror"
	"github.com/nightlyone/lockfile"
	log "github.com/sirupsen/logrus"
)

// Config defines the configuration interface
type Config interface {
	GetFirstRun() bool
	SetFirstRun(bool)
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
	log.WithField("path", path).
		Info("Loading config")

	if os.IsNotExist(err) {
		if flagUseConfig != "" {
			log.WithFields(log.Fields{
				"--config": flagUseConfig,
			}).Fatal("No user-provided configuration file was found")
		}
		log.Info(configFilename + " was not found. Generating one")
		c.SetFirstRun(true)
		if err = saveConfig(c, path, lockFile); err != nil {
			return
		}
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		c.SetFirstRun(true)
		if err = saveConfig(c, path, lockFile); err != nil {
			return
		}
	}

	// var jsonFile *os.File
	f, err := os.Open(path)
	onerror.Panic(err)
	defer f.Close()

	bv, _ := io.ReadAll(f)

	json.Unmarshal(bv, c)
	return
}

func saveConfig(c Config, path, lockFile string) (err error) {
	c.SetDefaults()
	log.WithField("path", path).
		Info("Saving config")

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
