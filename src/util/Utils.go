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

package util

import (
    "fmt"
    "log"
    "flag"
    "os"
    "io"
    "crypto/rand"
    "logger"
)

func Uuidgen() (string) {
        uuid := make([]byte, 16)
        n, err := io.ReadFull(rand.Reader, uuid)
        if n != len(uuid) || err != nil {
                //return "", err
                return ""
        }
        uuid[8] = uuid[8]&^0xc0 | 0x80
        uuid[6] = uuid[6]&^0xf0 | 0x40
        return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]) 
}

func Die(msg string) {
    logger.Log("verbose", msg)
    log.Printf("Usage:")
    flag.PrintDefaults()
    os.Exit(1)
}
