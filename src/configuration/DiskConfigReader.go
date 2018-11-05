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

package configuration

import (
    "fmt"
    "strings"
    "os"
    "os/exec"
    "util"
    "parameter"
    "logger"
)

type DiskConfigReader struct{
    targets, sources []string
    params parameter.Parameters
}
func (p *DiskConfigReader) Parameterize(params parameter.Parameters) {
    p.params = params
}
func (p *DiskConfigReader) SourceList() ([]string) {

    domain := p.params.Domain
    disks := strings.Split(p.params.Disks, ",")

    out, err := exec.Command("/usr/bin/virsh", "domblklist " + domain).Output()
    if err != nil {
	logger.Log("verbose", fmt.Sprintf("Error: domain %s: can't get list of disks\n %s", domain, err.Error()))
        os.Exit(1)
    }

    lines := strings.Split(string(out), "\n")
    has_disk := false

    for _, disk := range disks {
        has_disk = false
        for _, line := range lines {
            if strings.HasPrefix(line, disk) {
                target := line[:3]
                logger.Log("verbose", fmt.Sprintf("domain %s: detected target %s", domain, target))
                if target == disk {
                    p.targets = append(p.targets, target)
                    source := strings.Trim(line[4:len(line)], " ")
                    logger.Log("verbose", fmt.Sprintf("domain %s: detected source %s", domain, source))
                    p.sources = append(p.sources, source)
                    has_disk = true
		} 
            }
        }
        if ! has_disk {
	    logger.Log("verbose", fmt.Sprintf("Error: domain %s: disk %s not found", domain, disk))
            os.Exit(1)
        }
    }
    return p.sources
}
func (p *DiskConfigReader) TargetList() ([]string) {
    return p.targets
}
func (p *DiskConfigReader) SnapshotList() ([]string) {

    domain := p.params.Domain
    snapshots := []string{}

    for _, source := range p.sources {
        snapshot := source + "_" + util.Uuidgen()
        logger.Log("verbose", fmt.Sprintf("domain %s: using snapshot %s", domain, snapshot))
        snapshots = append(snapshots, snapshot)
    }
    return snapshots
}
func (p *DiskConfigReader) DestinationList() ([]string) {

    domain := p.params.Domain
    snapstore := p.params.Snapstore
    destinations := []string{}

    for _, target := range p.targets {
        destination := snapstore + "/" + domain + "_" + target 
        logger.Log("verbose", fmt.Sprintf("domain %s: using destination %s", domain, destination))
        destinations = append(destinations, destination)
    }
    return destinations
}
