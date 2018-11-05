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
)

const (
    msg_running_task = "domain %s: running task %s"
    msg_error_running_task = "domain %s: error while processing task %s"
    msg_task_not_found = "domain %s: task %s not found"
    msg_running_undo_operation = "domain %s: running undo operation for task %s"
)

type Executor struct{
}
func (e *Executor) Execute(params parameter.Parameters) bool {

    tasks := []string{"domain_suspend", "snapshot_create", "domain_resume", "snapshot_save", "sha256sum", "snapshot_remove",}

    if params.State == true {
        tasks = []string{"domain_save", "snapshot_create", "domain_restore", "snapshot_save", "sha256sum", "snapshot_remove",}
    }

    index := 0
    undo := false

    for i, task := range tasks {
        if ! execute(task,undo,params) {
            logger.Log(logger.Notice, fmt.Sprintf(msg_error_running_task, params.Domain, task));
	    index = i
	    undo = true
	    break
        }
    }
    if undo {
        for i := index; i >= 0; i-- {
	    task := tasks[i]
	    execute(task,undo,params)
	}
    }
    return true
}
func execute(name string, undo bool, params parameter.Parameters) bool {

    tasks := map[string]Task{
        "domain_suspend": &TaskVirshDomainSuspend{},
        "domain_resume": &TaskVirshDomainResume{},
        "domain_save": &TaskVirshDomainSave{},
        "domain_restore": &TaskVirshDomainRestore{},
        "snapshot_create": &TaskLVMSnapshotCreate{},
        "snapshot_save": &TaskLVMSnapshotSave{},
        "snapshot_remove": &TaskLVMSnapshotRemove{},
	"sha256sum": &TaskSHA256Sum{},
    }
    if task := tasks[name]; task == nil {
	logger.Log(logger.Notice, fmt.Sprintf(msg_task_not_found, params.Domain, name))
        return false
    } else {
        task.Parameterize(params)
	status := false
	if undo {
            logger.Log(logger.Debug, fmt.Sprintf(msg_running_undo_operation, params.Domain, name))
            status = task.Undo()
	} else {
            logger.Log(logger.Debug, fmt.Sprintf(msg_running_task, params.Domain, name))
            status = task.Execute()
	}
        return status
    }
}
