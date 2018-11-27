// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/liumingmin/goutils/safego"
)

var (
	Global Logger
)

func init() {
	// auto load config from default position
	Global = NewDefaultLogger(DEBUG)
	file, _ := exec.LookPath(os.Args[0])
	dir := filepath.Dir(file)
	if _, err := os.Stat("log4go.xml"); !os.IsNotExist(err) {
		Global.LoadConfiguration("log4go.xml")
	} else if _, err := os.Stat(filepath.Join(dir, "/log4go.xml")); !os.IsNotExist(err) {
		Global.LoadConfiguration(filepath.Join(dir, "log4go.xml"))
	} else if _, err := os.Stat(filepath.Join(dir, "/conf/log4go.xml")); !os.IsNotExist(err) {
		Global.LoadConfiguration(filepath.Join(dir, "/conf/log4go.xml"))
	} else {
		//fmt.Fprintf(os.Stderr, "log4go config not found, exec dir is: %s, u need to load it by yourself.\n", dir)
	}
}

// setup by config string, not a config file
func Setup(config []byte) {
	Global.Config(config)
}

// Wrapper for (*Logger).LoadConfiguration
func LoadConfiguration(filename string) {
	Global.LoadConfiguration(filename)
}

// Wrapper for (*Logger).AddFilter
func AddFilter(name string, lvl Level, writer LogWriter) {
	Global.AddFilter(name, lvl, writer)
}

// Wrapper for (*Logger).Close (closes and removes all logwriters)
func Close() {
	Global.Close()
}

func Crash(deep int, args ...interface{}) {
	if len(args) > 0 {
		Global.IntLogf(deep, CRITICAL, strings.Repeat(" %v", len(args))[1:], args...)
	}
	panic(args)
}

// Logs the given message and crashes the program
func Crashf(deep int, format string, args ...interface{}) {
	Global.IntLogf(deep, CRITICAL, format, args...)
	Global.Close() // so that hopefully the messages get logged
	panic(fmt.Sprintf(format, args...))
}

// Compatibility with `log`
func Exit(deep int, args ...interface{}) {
	if len(args) > 0 {
		Global.IntLogf(deep, ERROR, strings.Repeat(" %v", len(args))[1:], args...)
	}
	Global.Close() // so that hopefully the messages get logged
	os.Exit(0)
}

// Compatibility with `log`
func Exitf(deep int, format string, args ...interface{}) {
	Global.IntLogf(deep, ERROR, format, args...)
	Global.Close() // so that hopefully the messages get logged
	os.Exit(0)
}

// Compatibility with `log`
func Stderr(deep int, args ...interface{}) {
	if len(args) > 0 {
		Global.IntLogf(deep, ERROR, strings.Repeat(" %v", len(args))[1:], args...)
	}
}

// Compatibility with `log`
func Stderrf(deep int, format string, args ...interface{}) {
	Global.IntLogf(deep, ERROR, format, args...)
}

// Compatibility with `log`
func Stdout(deep int, args ...interface{}) {
	if len(args) > 0 {
		Global.IntLogf(deep, INFO, strings.Repeat(" %v", len(args))[1:], args...)
	}
}

// Compatibility with `log`
func Stdoutf(deep int, format string, args ...interface{}) {
	Global.IntLogf(deep, INFO, format, args...)
}

// Send a log message manually
// Wrapper for (*Logger).Log
func Log(lvl Level, source, message string) {
	Global.Log(lvl, source, message)
}

// Send a formatted log message easily
// Wrapper for (*Logger).Logf
func Logf(deep int, lvl Level, format string, args ...interface{}) {
	Global.IntLogf(deep, lvl, format, args...)
}

// Send a closure log message
// Wrapper for (*Logger).Logc
func Logc(deep int, lvl Level, closure func() string) {
	Global.IntLogf(deep, lvl, closure())
}

// Utility for finest log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Finest
func Finest(arg0 interface{}, args ...interface{}) {
	const (
		lvl = FINEST
	)
	if !IsFinestEnabled() {
		return
	}
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.IntLogf(DEFAULTL_DEEP, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.IntLogf(DEFAULTL_DEEP, lvl, first())
	default:
		// Build a format string so that it will be similar to Sprint
		Global.IntLogf(DEFAULTL_DEEP, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for fine log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Fine
func Fine(arg0 interface{}, args ...interface{}) {
	const (
		lvl = FINE
	)
	if !IsFineEnabled() {
		return
	}
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.IntLogf(DEFAULTL_DEEP, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.IntLogf(DEFAULTL_DEEP, lvl, first())
	default:
		// Build a format string so that it will be similar to Sprint
		Global.IntLogf(DEFAULTL_DEEP, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for debug log messages
// When given a string as the first argument, this behaves like Logf but with the DEBUG log level (e.g. the first argument is interpreted as a format for the latter arguments)
// When given a closure of type func()string, this logs the string returned by the closure iff it will be logged.  The closure runs at most one time.
// When given anything else, the log message will be each of the arguments formatted with %v and separated by spaces (ala Sprint).
// Wrapper for (*Logger).Debug
func Debug(arg0 interface{}, args ...interface{}) {
	const (
		lvl = DEBUG
	)
	if !IsDebugEnabled() {
		return
	}
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.IntLogf(DEFAULTL_DEEP, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.IntLogf(DEFAULTL_DEEP, lvl, first())
	default:
		// Build a format string so that it will be similar to Sprint
		Global.IntLogf(DEFAULTL_DEEP, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for trace log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Trace
func Trace(arg0 interface{}, args ...interface{}) {
	const (
		lvl = TRACE
	)
	if !IsTraceEnabled() {
		return
	}
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.IntLogf(DEFAULTL_DEEP, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.IntLogf(DEFAULTL_DEEP, lvl, first())
	default:
		// Build a format string so that it will be similar to Sprint
		Global.IntLogf(DEFAULTL_DEEP, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for info log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Info
func Info(arg0 interface{}, args ...interface{}) {
	const (
		lvl = INFO
	)
	if !IsInfoEnabled() {
		return
	}
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.IntLogf(DEFAULTL_DEEP, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.IntLogf(DEFAULTL_DEEP, lvl, first())
	default:
		// Build a format string so that it will be similar to Sprint
		Global.IntLogf(DEFAULTL_DEEP, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for Access log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Info
func Access(arg0 interface{}, args ...interface{}) {
	const (
		lvl = ACCESS
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.IntLogf(DEFAULTL_DEEP, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.IntLogf(DEFAULTL_DEEP, lvl, first())
	default:
		// Build a format string so that it will be similar to Sprint
		Global.IntLogf(DEFAULTL_DEEP, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for warn log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Warn
func Warn(arg0 interface{}, args ...interface{}) error {
	const (
		lvl = WARNING
	)
	if !IsWarnEnabled() {
		return nil
	}
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.IntLogf(DEFAULTL_DEEP, lvl, first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.IntLogf(DEFAULTL_DEEP, lvl, str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.IntLogf(DEFAULTL_DEEP, lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
	return nil
}

// Utility for error log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Error
func Error(arg0 interface{}, args ...interface{}) error {
	const (
		lvl = ERROR
	)
	if !IsErrorEnabled() {
		return nil
	}
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.IntLogf(DEFAULTL_DEEP, lvl, first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.IntLogf(DEFAULTL_DEEP, lvl, str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.IntLogf(DEFAULTL_DEEP, lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
	return nil
}

// Utility for critical log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Critical. This method will log the call stack
func Critical(arg0 interface{}, args ...interface{}) error {
	const (
		lvl = CRITICAL
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg := fmt.Sprintf("%s\n%s", fmt.Sprintf(first, args...), CallStack(DEFAULTL_DEEP))
		Global.IntLogf(DEFAULTL_DEEP, lvl, msg)
		//Global.IntLogf(lvl, CallStack(DEFAULTL_DEEP))
		return errors.New(fmt.Sprintf(first, args...))
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.IntLogf(DEFAULTL_DEEP, lvl, "%s\n%s", str, CallStack(DEFAULTL_DEEP))
		//Global.IntLogf(lvl, CallStack(DEFAULTL_DEEP))
		return errors.New(str)
	case func(interface{}) string:
		str := first(args[0])
		Global.IntLogf(DEFAULTL_DEEP, lvl, "%s\n%s", str, CallStack(DEFAULTL_DEEP))
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		msg := fmt.Sprintf("%s\n%s", fmt.Sprint(first)+fmt.Sprintf(strings.Repeat(" %v", len(args)), args...), CallStack(DEFAULTL_DEEP))
		Global.IntLogf(DEFAULTL_DEEP, lvl, msg)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
	return nil
}

// Recover used to log the stack when panic occur.
// usage: defer log4go.Recover("this is a msg: %v", "msg")
// or:
//      defer log4go.Recover(func(err interface{}) string {
//          // ... your code here, return the error message
//          return fmt.Sprintf("recover..v1=%v;v2=%v;err=%v", 1, 2, err)
//      })
func Recover(arg0 interface{}, args ...interface{}) {
	const (
		lvl = CRITICAL
	)
	if err := recover(); err != nil {
		switch first := arg0.(type) {
		case func(interface{}) string:
			// the recovered err will pass to this func
			//Critical(arg0, append([]interface{}{err}, args)...)
			Global.IntLogf(DEFAULTL_DEEP+2, lvl, "%s\n%s", first(err), CallStack(DEFAULTL_DEEP+2))
		case string:
			//Critical(a+"\n%v", append(args, err)...)
			msg := fmt.Sprintf("%s\n%s", fmt.Sprintf(first, args...), CallStack(DEFAULTL_DEEP))
			Global.IntLogf(DEFAULTL_DEEP, lvl, msg)
		default:
			//Critical(arg0, append(args, err)...)
			msg := fmt.Sprintf("%s\n%s", fmt.Sprint(first)+fmt.Sprintf(strings.Repeat(" %v", len(args)), args...), CallStack(DEFAULTL_DEEP))
			Global.IntLogf(DEFAULTL_DEEP, lvl, msg)
		}
	}
}

func IsFinestEnabled() bool {
	return isLevelEnabled(FINEST)
}

func IsFineEnabled() bool {
	return isLevelEnabled(FINE)
}

func IsDebugEnabled() bool {
	return isLevelEnabled(DEBUG)
}

func IsTraceEnabled() bool {
	return isLevelEnabled(TRACE)
}

func IsInfoEnabled() bool {
	return isLevelEnabled(INFO)
}

func IsWarnEnabled() bool {
	return isLevelEnabled(WARNING)
}

func IsErrorEnabled() bool {
	return isLevelEnabled(ERROR)
}

func isLevelEnabled(lvl Level) bool {
	//fmt.Printf("++++++++  lowest level: %v, curr Level: %v, isEnable: %v \n", LowestLevel, lvl, (lvl >= LowestLevel))
	return lvl >= LowestLevel
	//for _, filt := range Global {
	//if lvl >= filt.Level {
	//// return true if any filt matched
	//return true
	//}
	//}
	//return false
}
