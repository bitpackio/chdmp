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
    "strings"
    "os/exec"
    "logger"
    "parameter"
)

const (
    msg_computing_snapshot_hash = "domain %s: computing hash for snapshot %s"
    msg_error_computing_snapshot_hash = "domain %s : error while computing hash for snapshot %s"
    msg_computing_destination_hash = "domain %s: computing hash for destination %s"
    msg_error_computing_destination_hash = "domain %s : error while computing hash for destination %s"
    msg_hash = "domain %s: hash: %s"
    msg_hashs_not_equal = "domain %s: hashs not equal"
    msg_hashs_equal = "domain %s: hashs for snapshot and destination verified and equal"
)

type TaskSHA256Sum struct{
    pid int
    params parameter.Parameters
}
func (p *TaskSHA256Sum) Execute() bool {

    return true

    domain_name := p.params.Domain

    snapshotHashes := []string{}

    for i := range p.params.Snapshots {
        snapshot := p.params.Snapshots[i]
        logger.Log(logger.Notice, fmt.Sprintf(msg_computing_snapshot_hash, domain_name, snapshot))
	out, err := exec.Command("/usr/bin/sha256sum", snapshot).Output()
        outFields := strings.Fields(fmt.Sprintf("%s",out))
	snapshotHash := outFields[0]
	snapshotHashes = append(snapshotHashes,snapshotHash)
	if err != nil {
            logger.Log(logger.Notice, fmt.Sprintf(msg_error_computing_snapshot_hash, domain_name, snapshot))
            return false
        }
    }

    destinationHashes := []string{}

    for i := range p.params.Destinations {
        destination := p.params.Destinations[i]
        logger.Log(logger.Notice, fmt.Sprintf(msg_computing_destination_hash, domain_name, destination))
	out, err := exec.Command("/usr/bin/sha256sum", destination).Output()
        outFields := strings.Fields(fmt.Sprintf("%s",out))
	destinationHash := outFields[0]
	destinationHashes = append(destinationHashes,destinationHash)
	if err != nil {
            logger.Log(logger.Notice, fmt.Sprintf(msg_error_computing_destination_hash, domain_name, destination))
            return false
        }
    }

    for i := range p.params.Snapshots {
        snapshotHash := snapshotHashes[i]
        destinationHash := destinationHashes[i]
	logger.Log(logger.Notice, fmt.Sprintf(msg_hash, domain_name, snapshotHash))
	logger.Log(logger.Notice, fmt.Sprintf(msg_hash, domain_name, destinationHash))
	if snapshotHash != destinationHash {
            logger.Log(logger.Notice, fmt.Sprintf(msg_hashs_not_equal, domain_name))
            return false
	}
        logger.Log(logger.Notice, fmt.Sprintf(msg_hashs_equal, domain_name))
    }
    return true
}
func (p *TaskSHA256Sum) Undo() bool {
    return true
}
func (p *TaskSHA256Sum) Parameterize(params parameter.Parameters) {
    p.params = params
}
