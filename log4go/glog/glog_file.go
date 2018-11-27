// Go support for leveled logs, analogous to https://code.google.com/p/google-glog/
//
// Copyright 2013 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// File I/O for logs.

package glog

import (
	// "errors"
	// "flag"
	"fmt"
	"os"
	// "path/filepath"
	// "sync"
	"time"
)

// MaxSize is the maximum size of a log file in bytes.
var MaxSize uint64 = 1024 * 1024 * 1800
var LogNameOfLevel = []string{"acce", "fnst", "fine", "debg", "trac", "info", "warn", "eror", "crit", "fatl"}

// logName returns a new log file name containing tag, with start time t, and
// the name for the symlink for tag.
func logName(level int32, t time.Time) (fname, rotateName string) {
	fname = LogNameOfLevel[level]
	rotateName = fmt.Sprintf("%s.%04d%02d%02d%02d%02d%02d",
		fname,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second())
	return fname, rotateName
}

// create creates a new log file and returns the file and its filename, which
// contains tag ("INFO", "FATAL", etc.) and t.  If the file is created
// successfully, create also attempts to update the symlink for that tag, ignoring
// errors.
func create(level int32, t time.Time) (f *os.File, filename string, err error) {
	fname, rotateName := logName(level, t)
	os.Rename(fname, rotateName) // ignore err
	f, err = os.Create(fname)
	if err == nil {
		return f, fname, nil
	}
	return nil, "", fmt.Errorf("log: cannot create log: %v", err)
}
