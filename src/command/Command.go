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

package command

import (
	"strings"
	"os/exec"
	"fmt"
	"logger"
	"parameter"
)

type Command struct {
    Params parameter.Parameters
    Cmd string
    Args []string
}

func (c *Command) Exec() bool {

    domain := c.Params.Domain

    logger.Log("debug", fmt.Sprintf("domain %s: " + c.Cmd + " " + strings.Join(c.Args, " "), domain))

    cmd := exec.Command(c.Cmd, c.Args...)

    err := cmd.Start()
    if err != nil {
        //logger.StdLog.Fatal(err)
        logger.Log("debug", fmt.Sprintf("%v", err))
        return false
    }

    pid := cmd.Process.Pid

    logger.Log("debug", fmt.Sprintf("domain %s: waiting for child process with pid %d to finish...", domain, pid))
    err = cmd.Wait()
    if err != nil {
        logger.Log("debug", fmt.Sprintf("domain %s: child process exited with %v", domain, err))
        return false
    }
    return true
}
