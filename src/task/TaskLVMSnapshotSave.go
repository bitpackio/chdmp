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

package task

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "logger"
    "parameter"
    "stream"
)

const (
    msg_updating_chunks = "domain %s: updating chunks from %s to %s"
    msg_chunks_updated = "domain %s: chunks from %s to %s updated"
)

type TaskLVMSnapshotSave struct{
    pid int
    params parameter.Parameters
}
func (p *TaskLVMSnapshotSave) Execute() bool {

    dump := p.params.Dump
    domain_name := p.params.Domain

    for i := range p.params.Snapshots {

        target := p.params.Targets[i]
        snapshot := p.params.Snapshots[i]
        destination := p.params.Destinations[i]

        logger.Log(logger.Notice, fmt.Sprintf(msg_updating_chunks, domain_name, target, destination))

        if ! dump {
            continue
        }

        chunkStream := stream.ChunkStream{}

        signal_chan := make(chan os.Signal, 1)
        signal.Notify(signal_chan,
                syscall.SIGUSR1,
                syscall.SIGHUP,
                syscall.SIGINT,
                syscall.SIGTERM,
                syscall.SIGQUIT)
        go chunkStream.ProcessSignals(signal_chan)

        p.params.Src = snapshot
        p.params.Dst = destination

        chunkStream.Parameterize(p.params)
        chunkStream.Initialize()
        chunkStream.Execute()
        chunkStream.Statistics()

        logger.Log(logger.Verbose, fmt.Sprintf(msg_chunks_updated, domain_name, target, destination))
    }
    return true
}
func (p *TaskLVMSnapshotSave) Undo() bool {
    return true
}
func (p *TaskLVMSnapshotSave) Parameterize(params parameter.Parameters) {
    p.params = params
}
