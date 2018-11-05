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
    "strconv"
    "os"
    "os/exec"
    "errors"
    "logger"
    "parameter"
    "command"
)

const (
    msg_creating_snapshot = "domain %s: creating snapshot %s"
    msg_error_creating_snapshot = "domain %s : error while creating snapshot %s"
    msg_snapshot_created = "domain %s: snapshot %s created"
)

type TaskLVMSnapshotCreate struct{
    pid int
    params parameter.Parameters
}
func (p *TaskLVMSnapshotCreate) Execute() bool {

    domain_name := p.params.Domain
    snapsize := strconv.Itoa(p.params.Snapsize) + "G"

    for i := range p.params.Snapshots {

        snapshot := p.params.Snapshots[i]
        source := p.params.Sources[i]

        logger.Log(logger.Debug, fmt.Sprintf("verifying source: %s", source))

	err := p.Verify(source)
	if err != nil {
	    logger.Log(logger.Notice, err.Error())
	    return false
        }

        cmd_name := "/bin/sh"
        cmd_args := []string{
            "-c",
            "/sbin/lvcreate --chunksize 512 --permission r --snapshot --size " + snapsize + " --name " + snapshot + " " + source,
        }
        cmd := command.Command{p.params, cmd_name, cmd_args}

        logger.Log(logger.Notice, fmt.Sprintf(msg_creating_snapshot, domain_name, snapshot))
        if err := cmd.Exec(); err != true {
            logger.Log(logger.Notice, fmt.Sprintf(msg_error_creating_snapshot, domain_name, snapshot))
            return false
        }
        logger.Log(logger.Debug, fmt.Sprintf(msg_snapshot_created, domain_name, snapshot))
    }
    return true
}
func (p *TaskLVMSnapshotCreate) Undo() bool {
    snapshotRemove := TaskLVMSnapshotRemove{}
    snapshotRemove.Parameterize(p.params)
    return snapshotRemove.Execute()
}
func (p *TaskLVMSnapshotCreate) Parameterize(params parameter.Parameters) {
    p.params = params
}
func (p *TaskLVMSnapshotCreate) Verify(source string) error {

    domain := p.params.Domain

    if _, err := os.Stat(source); os.IsNotExist(err) {
	logger.Log(logger.Debug, fmt.Sprintf("domain %s: Error: %s", domain, err.Error()))
        return errors.New(fmt.Sprintf("domain %s: Error: source %s does not exists", domain, source))
    }

    vg := ""
    lv := ""
    if strings.HasPrefix(source, "/dev/") {
        vg = source[5:9]
        lv = source[10:len(source)]
        logger.Log(logger.Debug, fmt.Sprintf("domain %s: detected volume group %s", domain, vg))
        logger.Log(logger.Debug, fmt.Sprintf("domain %s: detected logical volume %s", domain, lv))
    } else {
	return errors.New(fmt.Sprintf("domain %s: Error: can't parse volume group and logical volume", domain))
    }
    out, err := exec.Command("/sbin/vgdisplay", vg).Output()
    if err != nil {
        logger.Log(logger.Debug, (fmt.Sprintln("domain %s: Error: %s", domain, err.Error())))
        return errors.New(fmt.Sprintln("domain %s: Error: no volume group found\n %s", domain))
    }
    vg_lines := strings.Split(string(out), "\n")
    for _, vg_line := range vg_lines {
        logger.Log(logger.Debug, fmt.Sprintf("domain %s: parsing lvm data: %s", domain, vg_line))
        if strings.HasPrefix(vg_line, "  VG Size") {
            vg_raw := strings.Trim(vg_line[11:len(vg_line)], " ")
            i := strings.IndexByte(vg_raw, ' ')
            if i == -1 {
                return errors.New(fmt.Sprintf("domain %s: Error: can't parse vg_raw with value %s", domain, vg_raw))
            }
            vg_size := vg_raw[0:i]
            vg_unit := strings.Trim(vg_raw[i:len(vg_raw)], " ")
            logger.Log(logger.Debug, fmt.Sprintf("domain %s: detected vg_size %s with unit %s", domain, vg_size, vg_unit))
        }
        if strings.HasPrefix(vg_line, "  Free  PE / Size") {
            vg_raw := strings.Trim(vg_line[18:len(vg_line)], " ")
            i := strings.IndexByte(vg_raw, '/')
            if i != -1 {
                vg_raw = vg_raw[i+2:len(vg_raw)] 
            } else {
                return errors.New(fmt.Sprintf("domain %s: Error: can't parse vg_raw with value %s", domain, vg_raw))
            }
	    i = strings.IndexByte(vg_raw, ' ')
            if i == -1 {
                return errors.New(fmt.Sprintf("domain %s: Error: can't parse vg_raw with value %s", domain, vg_raw))
            }
            vg_free_pe := vg_raw[0:i]
            vg_unit := strings.Trim(vg_raw[i:len(vg_raw)], " ")
            logger.Log(logger.Debug, fmt.Sprintf("domain %s: detected vg_free_pe %s with unit %s", domain, vg_free_pe, vg_unit))

            i = strings.IndexByte(vg_free_pe, ',')
            if i == -1 {
                i = strings.IndexByte(vg_free_pe, '.')
                if i == -1 {
                    return errors.New(fmt.Sprintf("domain %s: Error: can't parse vg_free_pe with value %s", domain, vg_free_pe))
                }
            }
            vg_free_pe_rounded := vg_free_pe[0:i]
            vg_free_pe_int, err := strconv.Atoi(vg_free_pe_rounded)
            if err != nil {
                logger.Log(logger.Debug, fmt.Sprintf("domain %s: Error: %s", err.Error()))
                return errors.New(fmt.Sprintf("domain %s: Error: can't convert vg free pe %s", vg_free_pe))
            }
            if p.params.Snapsize > vg_free_pe_int {
                return errors.New(fmt.Sprintf("domain %s: Error: volume group does not have enough free space to hold snapshots", domain))
            } else {
                logger.Log(logger.Debug, fmt.Sprintf("domain %s: detected vg_free_pe with %s %s free space can hold snapshots with size of %s %s", 
                            domain, strconv.Itoa(vg_free_pe_int), vg_unit, strconv.Itoa(p.params.Snapsize), vg_unit))
            }
        }
        if strings.HasPrefix(vg_line, "  VG UUID") {
            vg_uuid := strings.Trim(vg_line[11:len(vg_line)], " ")
            logger.Log(logger.Debug, fmt.Sprintf("domain %s: detected vg_uuid %s", domain, vg_uuid))
        }
    }
    return nil
}
