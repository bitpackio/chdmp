/*
  chdmp - chunk based stream dump

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

package main

import (
    "flag"
    "log"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "parameter"
    "stream"
)

func Usage() {
    fmt.Println("Usage:")
    flag.PrintDefaults()
    os.Exit(0)
}

func main() {

    chunkStream := stream.ChunkStream{}

    signal_chan := make(chan os.Signal, 1)
    signal.Notify(signal_chan,
                syscall.SIGUSR1,
                syscall.SIGHUP,
                syscall.SIGINT,
                syscall.SIGTERM,
                syscall.SIGQUIT)
    go chunkStream.ProcessSignals(signal_chan)

    stats := flag.Bool("stats", true, "enable or disable statistic output")
    debug := flag.Bool("debug", false, "enable or disable debug logging")
    verbose := flag.Bool("verbose", false, "enable or disable verbose logging")
    hash := flag.Bool("hash", false, "enable or disable sha256 hash mode")
    simulate := flag.Bool("simulate", false, "enable or disable simulation mode")
    force := flag.Bool("force", false, "enable or disable force mode")
    uselog := flag.Bool("log", false, "enable or disable log mode")
    flush := flag.Bool("flush", false, "enable or disable flush to disk mode")
    help := flag.Bool("help", false, "help usage")
    chunkSize := flag.Int("chunksize", 8192, "set the chunk-size to be used")
    var src string
    flag.StringVar(&src, "input", "", "specify the input stream")
    var dst string
    flag.StringVar(&dst, "output", "", "specify the output stream")

    flag.Parse()
 
    if src == ""  {
        log.Printf("Error: option input does not contains a value")
        Usage()
    } 
    if dst == ""  {
        log.Printf("Error: option output does not contains a value")
        Usage()
    } 
    if *help {
        Usage()
    } 
    params := parameter.Parameters{}
    params.Stats = *stats
    params.Debug = *debug
    params.Verbose = *verbose
    params.Hash = *hash
    params.Simulate = *simulate
    params.Force = *force
    params.Uselog = *uselog
    params.Flush = *flush
    params.ChunkSize = *chunkSize
    params.Src = src
    params.Dst = dst

    chunkStream.Parameterize(params)
    chunkStream.Initialize()
    chunkStream.Execute()
    chunkStream.Statistics()
}
