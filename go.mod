module github.com/jwmwalrus/rtcycler

go 1.24.0

toolchain go1.24.4

require (
	github.com/adrg/xdg v0.5.3
	github.com/jwmwalrus/bnp v1.23.1
	github.com/nightlyone/lockfile v1.0.0
	github.com/pborman/getopt/v2 v2.1.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require golang.org/x/sys v0.39.0 // indirect

// replace github.com/jwmwalrus/bnp => ../bnp
