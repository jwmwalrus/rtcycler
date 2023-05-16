package rtc

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/jwmwalrus/bnp/ing2"
	"github.com/pborman/getopt/v2"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	varBasePath = "/var/lib"
)

var (
	appDirName      string
	appInstance     string
	appName         = "app"
	cacheDir        string
	conf            Config
	configDir       string
	configFile      string
	configFilename  = "config.json"
	daemonDir       string
	dataDir         string
	flagDaemonMode  bool
	flagDebug       bool
	flagDry         bool
	flagEchoLogging bool
	flagHelp        bool
	flagLogLevel    string
	flagTestMode    bool
	flagUseConfig   string
	flagVerbose     bool
	instanceSuffix  string
	lockFile        string
	lockFilename    string
	logFile         *lumberjack.Logger
	logFilename     = "app.log"
	runtimeDir      string
)

func init() {
	getopt.FlagLong(&flagHelp, "help", 'h', "Display this help")
	getopt.FlagLong(&flagDry, "dry-run", 'n', "Dry run")
	getopt.FlagLong(&flagDaemonMode, "daemon", 0, "Daemon mode. Uses /var/lib as XDG base (Unix-like only)")
	getopt.FlagLong(&flagDebug, "debug", 0, "Start with logging debug level")
	getopt.FlagLong(&flagTestMode, "test", 0, "Start in test mode")
	getopt.FlagLong(&flagVerbose, "verbose", 'v', "Bump logging level")
	getopt.FlagLong(&flagLogLevel, "log-level", 0, "Logging level (debug|info|warn|error|fatal)")
	getopt.FlagLong(&flagEchoLogging, "echo-logging", 'e', "Echo logs to stderr")
	getopt.FlagLong(&flagUseConfig, "config", 'c', "Use provided config file")
}

// AppDirName returns the passed application's directory name
func AppDirName() string { return appDirName }

// AppInstance returns the current application's instance
func AppInstance() string { return appInstance }

// AppName returns the passed application's name
func AppName() string { return appName }

// CacheDir returns the XDG's home directory for cache
func CacheDir() string { return cacheDir }

// ConfigDir returns the XDG's home directory for config
func ConfigDir() string { return configDir }

// ConfigFilename returns the passed config filename
func ConfigFilename() string { return configFilename }

// DaemonDir returns the /var/dir subdirectory for the application
func DaemonDir() string { return daemonDir }

// DataDir returns the XDG's home directory for data
func DataDir() string { return dataDir }

// FlagDaemonMode returns the value of the --daemon flag
func FlagDaemonMode() bool { return flagDaemonMode }

// FlagDebug returns the value of the --debug flag
func FlagDebug() bool { return flagDebug }

// FlagDry returns the value of the --dry-run flag
func FlagDry() bool { return flagDry }

// FlagEchoLogging returns the value of the --echo-logging flag
func FlagEchoLogging() bool { return flagEchoLogging }

// FlagLogLogLevel returns the value of the --log-level flag
func FlagLogLogLevel() string { return flagLogLevel }

// FlagTestMode returns the value of the --test-mode flag
func FlagTestMode() bool { return flagTestMode }

// FlagUseConfig returns the value of the --config flag
func FlagUseConfig() string { return flagUseConfig }

// FlagVerbose returns the value of the --verbose flag
func FlagVerbose() bool { return flagVerbose }

// InstanceConfig instance configuration
func InstanceConfig() Config { return conf }

// InstanceSuffix suffix used for the running instance
func InstanceSuffix() string {
	if instanceSuffix == "" {
		instanceSuffix, _ = ing2.GetRandomLetters(8)
	}
	return instanceSuffix
}

// LockFilename returns the passed lock filename
func LockFilename() string { return lockFilename }

// LogFilename returns the passed log filename
func LogFilename() string { return logFilename }

// OS returns the current OS
func OS() string { return runtime.GOOS }

// ResetInstanceSuffix clears the instanceSuffix value.
func ResetInstanceSuffix() { instanceSuffix = "" }

// RuntimeDir returns XDG's run (volatile) directory
func RuntimeDir() string { return runtimeDir }

// SetTestMode sets the value of the --test-mode flag
func SetTestMode() {
	flagTestMode = true
}

// UnsetTestMode unsets the value of the --test-mode flag
func UnsetTestMode() {
	flagTestMode = false
}

func parseArgs() (args []string) {
	getopt.Parse()
	args = getopt.Args()
	arg0 := []string{os.Args[0]}
	args = append(arg0, args...)

	if flagHelp {
		getopt.Usage()
		os.Exit(0)
	}

	resolveLogLevel()

	if flagEchoLogging {
		mw := io.MultiWriter(os.Stderr, logFile)
		log.SetOutput(mw)
	}

	return
}

func resolveLogLevel() {
	givenLogLevel := flagLogLevel

	if givenLogLevel == "" {
		if flagDebug {
			flagLogLevel = "debug"
		} else if flagTestMode {
			flagLogLevel = "debug"
		} else if flagVerbose {
			flagLogLevel = "info"
		} else {
			flagLogLevel = "error"
		}
	} else {
		if _, err := log.ParseLevel(givenLogLevel); err != nil {
			fmt.Printf("Unsupported log level: %v\n", givenLogLevel)
			flagLogLevel = "error"
		} else {
			flagLogLevel = givenLogLevel
		}
	}

	level, _ := log.ParseLevel(flagLogLevel)
	log.SetLevel(level)
	log.SetReportCaller(flagLogLevel == "debug")
}
