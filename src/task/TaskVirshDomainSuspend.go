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
    msg_suspending = "domain %s: suspending"
    msg_error_suspending = "domain %s: suspending"
    msg_suspended = "domain %s: suspended"
)

type TaskVirshDomainSuspend struct {
    pid int
    params parameter.Parameters
}
func (p *TaskVirshDomainSuspend) Execute() bool {

    domain_name := p.params.Domain

    cmd_name := "virsh"
    cmd_args := []string{
        "suspend",
        domain_name,
    }

    cmd := command.Command{p.params, cmd_name, cmd_args}

    logger.Log(logger.Notice, fmt.Sprintf(msg_suspending, domain_name))
    if err := cmd.Exec(); err != true {
        logger.Log(logger.Notice, fmt.Sprintf(msg_error_suspending, domain_name))
        return false
    }
    logger.Log(logger.Verbose, fmt.Sprintf(msg_suspended, domain_name))
    return true
}
func (p *TaskVirshDomainSuspend) Undo() bool {
    domainResume := TaskVirshDomainResume{}
    domainResume.Parameterize(p.params)
    return domainResume.Execute()
}
func (p *TaskVirshDomainSuspend) Parameterize(params parameter.Parameters) {
    p.params = params
}
