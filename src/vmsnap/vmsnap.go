/*
  vmsnap - save kvm domain which uses logical volume block devices 

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
    "fmt"
    "os"
    "flag"
    "util"
    "logger"
    "parameter"
    "configuration"
    "task"
)

func init() {
    logger.Stdout = true
    logger.Syslog = true
    logger.Filelog = false
    logger.Logpath =  "/var/log/vmsnap.log"
}

func main() {

    debug := flag.Bool("debug", false, "enable or disable debug logging")
    verbose := flag.Bool("verbose", false, "enable or disable verbose logging")
    dump := flag.Bool("dump", false, "enable or disable disk dump to file")
    state := flag.Bool("state", false, "enable or disable writing domain memory state to file")
    snapsize := flag.Int("snapsize", 2, "set the logical volume size between 1 to 30 in gigabytes for snapshots")
    var snapstore string
    flag.StringVar(&snapstore, "snapstore", "/var/lib/libvirt/backups", "set the directory to write domain states and disk dumps")
    var disks string
    flag.StringVar(&disks, "disks", "", "specify the disks to dump")
    var domain string
    flag.StringVar(&domain, "domain", "", "specify the domain to process")

    stats := flag.Bool("stats", true, "enable or disable statistic output")
    hash := flag.Bool("hash", false, "enable or disable sha256 hash mode")
    simulate := flag.Bool("simulate", false, "enable or disable simulation mode")
    force := flag.Bool("force", false, "enable or disable force mode")
    flush := flag.Bool("flush", false, "enable or disable flush to disk mode")
    chunkSize := flag.Int("chunksize", 4098, "set the chunk-size to be used")
    //help := false //flag.Bool("help", false, "help usage")

    flag.Parse()

    if *snapsize < 1 || *snapsize > 20  {
        util.Die(fmt.Sprintf("Error: given snapsize value %d not valid", *snapsize))
    } 
    if snapstore == ""  {
        util.Die(fmt.Sprintf("Error: option snapstore does not contains a value"))
    } else if snapstore == "/" {
        util.Die(fmt.Sprintf("Error: enabling backup directory / does not allowed"))
    }
    if _, err := os.Stat(snapstore); os.IsNotExist(err) {
        util.Die(fmt.Sprintf("Error: snapstore %s does not exist ", snapstore))
    }
    if domain == ""  {
        util.Die(fmt.Sprintf("Error: option domain does not contains a value"))
    } 
    if _, err := os.Stat("/etc/libvirt/qemu/" + domain + ".xml"); os.IsNotExist(err) {
        util.Die(fmt.Sprintf("Error: domain %s not defined", domain))
    }
    if disks == ""  {
        util.Die(fmt.Sprintf("Error: option disks does not contains a value"))
    } 

    params := parameter.Parameters{}
    params.Debug = *debug
    params.Verbose = *verbose
    params.Dump = *dump
    params.State = *state
    params.Snapsize = *snapsize
    params.Snapstore = snapstore
    params.Domain = domain
    params.Disks = disks

    params.Stats = *stats
    params.Hash = *hash
    params.Simulate = *simulate
    params.Force = *force
    params.Flush = *flush
    params.ChunkSize = *chunkSize

    logger.Log("debug", fmt.Sprintf("domain %s: pid %d", domain, os.Getpid()))

    configuration := configuration.DiskConfigReader{}
    configuration.Parameterize(params)
    params.Sources = configuration.SourceList()
    params.Targets = configuration.TargetList()
    params.Snapshots = configuration.SnapshotList()
    params.Destinations = configuration.DestinationList()

    executor := task.Executor{}
    executor.Execute(params)
}
