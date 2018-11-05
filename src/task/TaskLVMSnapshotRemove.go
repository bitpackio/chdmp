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
    "errors"
    "logger"
    "parameter"
    "command"
)

const (
    msg_removing_snapshot = "domain %s: removing snapshot %s"
    msg_error_removing_snapshot = "domain %s : error while removing snapshot %s"
    msg_snapshot_removed = "domain %s: snapshot %s removed"
    msg_verify_snapshot = "domain %s: error: %s"
)

type TaskLVMSnapshotRemove struct{
    pid int
    params parameter.Parameters
}
func (p *TaskLVMSnapshotRemove) Execute() bool {

    domain_name := p.params.Domain

    for i := range p.params.Snapshots {

        snapshot := p.params.Snapshots[i]

	err := p.Verify(snapshot)
	if err != nil {
	    logger.Log("debug", err.Error())
	    return false
        }

        cmd_name := "/bin/sh"
        cmd_args := []string{
            "-c",
            "/sbin/lvremove --force " + snapshot,
        }

        cmd := command.Command{p.params, cmd_name, cmd_args}

        logger.Log(logger.Notice, fmt.Sprintf(msg_removing_snapshot, domain_name, snapshot))
        if err := cmd.Exec(); err != true {
            logger.Log(logger.Notice, fmt.Sprintf(msg_error_removing_snapshot, domain_name, snapshot))
            return false
        }
        logger.Log(logger.Debug, fmt.Sprintf(msg_snapshot_removed, domain_name, snapshot))
    }
    return true
}
func (p *TaskLVMSnapshotRemove) Undo() bool {
    return true
}
func (p *TaskLVMSnapshotRemove) Parameterize(params parameter.Parameters) {
    p.params = params
}
func (p *TaskLVMSnapshotRemove) Verify(snapshot string) error {
    if _, err := os.Stat(snapshot); os.IsNotExist(err) {
        return errors.New(fmt.Sprintf(msg_verify_snapshot, p.params.Domain, err.Error()))
    }
    return nil
}
