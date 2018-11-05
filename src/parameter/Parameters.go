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

package parameter

type Parameters struct{
    Debug, Verbose, Dump, State, Uselog bool
    Snapstore, Domain, Disks string
    Targets, Sources, Snapshots, Destinations []string
    Snapsize int

    Stats, Hash, Simulate, Force, Flush bool
    ChunkSize int
    Src, Dst  string
}
