/*
  Copyright (C) 2018 bitpack.io <hello@bitpack.io>

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License at <http://www.gnu.org/licenses/> for
  more details.
*/

package logger

import (
    "log"
    "log/syslog"
    "os"
)

const (
    Notice = "notice"
    Verbose = "verbose"
    Debug = "debug"
)

var (
  StdLog, SysLog,  FileLog *log.Logger
  Stdout , Syslog, Filelog bool
  Logpath string
  _Notice, _Verbose, _Debug bool
)

func init() {
    Stdout = true
    Syslog = false
    Filelog = false
    Logpath = "/var/log/vmsnap.log"
    _Notice = true
    _Verbose = true
    _Debug = false
    StdLog = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

func initSyslog() {
    Syslog = true
    logwriter, e := syslog.New(syslog.LOG_NOTICE, "vmsnap")
    if e == nil {
        SysLog = log.New(logwriter, "", log.Lshortfile)
    }
}

func initFilelog() {
    Filelog = true
    var file, err = os.Create(Logpath)
    if err != nil {
        panic(err)
    }
    FileLog = log.New(file, "", log.LstdFlags|log.Lshortfile)
}

func _log(msg string) {
    if Stdout {
        StdLog.Print(msg)
    }
    if Syslog {
	if SysLog == nil {
            initSyslog()
	}
        SysLog.Print(msg)
    }
    if Filelog {
	if FileLog == nil {
            initFilelog()
	}
        FileLog.Print(msg)
    }
}

func Log(level, msg string) {
    if level == "notice" {
        if(_Notice) {
            _log(msg)
        }
    }
    if level == "verbose" {
        if(_Verbose) {
            _log(msg)
        }
    }
    if level == "debug" {
        if(_Debug) {
            _log(msg)
        }
    }
}
