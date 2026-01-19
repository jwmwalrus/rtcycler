ChangeLog
=========

All noticeable changes in the project  are documented in this file.

Format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

This project uses [semantic versions](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.7.4] 2026-01-19

Update dependencies

## [0.7.3] 2026-01-06

Update dependencies

## [0.7.2] 2025-04-06

Update dependencies

## [0.7.1] 2024-12-21

Update dependencies

## [0.7.0] 2024-02-01

Granular changes

### Added

* `ConfigLocker` interface to make `SetLockFile`'s implementation optional
* `ConfigResolver` interface to be invoked after loading the config

### Modified

* The `Load` function now uses functional options
* The `Config` interface no longer requires `SetLockFile`

### Removed

## [0.6.6] 2024-01-10

Update dependencies

## [0.6.5] 2023-10-17

Update dependencies

## [0.6.4] 2023-09-09

Adjust to bnp change

## [0.6.3] 2023-09-02

Do not invoke dummy unloaders

## [0.6.2] 2023-08-29

Fix logger's LogLevel issue

## [0.6.1] 2023-08-28

Fix logger's reported source

## [0.6.0] 2023-08-28

Implement RuntimeLogger

## [0.5.0] 2023-08-24

Rework logging in terms of slog

## [0.4.0] 2023-08-23

Set Config's lockFile right after loading

## [0.3.0] 2023-08-18

Add WithDotHome option

## [0.2.1] 2023-08-13

Upgrade Go version

## [0.2.0] 2023-05-16

Add ResetInstanceSuffix function

## [0.1.2] 2023-05-15

 Fix instance suffix for testing

## [0.1.1] 2023-04-16

 Fix runtime dir issue when in daemon mode

## [0.1.0] 2023-04-01

Initial release
