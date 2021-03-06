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
    "logger"
    "parameter"
    "command"
)

const (
    msg_restoring_state = "domain %s: restoring state from %s"
    msg_error_restoring_state = "domain %s: error while restoring state from %s"
    msg_state_restored = "domain %s: state restored from %s"
)

type TaskVirshDomainRestore struct{
    pid int
    params parameter.Parameters
}
func (p *TaskVirshDomainRestore) Execute() bool {

    domain_name := p.params.Domain
    state_file := p.params.Snapstore + "/" + domain_name + ".state"

    cmd_name := "virsh"
    cmd_args := []string{
        "restore",
        state_file,
    }

    cmd := command.Command{p.params, cmd_name, cmd_args}

    logger.Log(logger.Notice, fmt.Sprintf(msg_restoring_state, domain_name, state_file))
    if err := cmd.Exec(); err != true {
        logger.Log(logger.Notice, fmt.Sprintf(msg_error_restoring_state, domain_name, state_file))
        return false
    }
    logger.Log(logger.Verbose, fmt.Sprintf(msg_state_restored, domain_name, state_file))
    return true
}
func (p *TaskVirshDomainRestore) Undo() bool {
    return true
}
func (p *TaskVirshDomainRestore) Parameterize(params parameter.Parameters) {
    p.params = params
}
