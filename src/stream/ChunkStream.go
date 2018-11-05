/*
  ChunkStream - chunk based stream processing

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

package stream

import (
    "log"
    "fmt"
    "io"
    "os"
    "bytes"
    "math"
    "syscall"
    "time"
    "crypto/sha256"
    "encoding/base64"
    "logger"
    "parameter"
)

type ChunkStream struct{
    params parameter.Parameters
    Src, Dst string 
    Stats, Debug, Verbose, Hash, Simulate, Force, Uselog, Flush bool
    ChunkSize int

    chunksIn, chunksOut int64 
    bytesIn, bytesOut int64 
    eofOnSrc, eofOnDest bool
    f1, f2, f3  *os.File 
    timeStart time.Time
}

func (p *ChunkStream) Parameterize(params parameter.Parameters) {
    p.params = params
    p.Src = p.params.Src
    p.Dst = p.params.Dst
    p.Debug = p.params.Debug
    p.Verbose = p.params.Verbose
    p.Hash = p.params.Hash
    p.Simulate = p.params.Simulate
    p.Force = p.params.Force
    p.Flush = p.params.Flush
    p.ChunkSize = p.params.ChunkSize
    p.timeStart = time.Now()
}

func (p *ChunkStream) Initialize() {
    p.Log("debug", fmt.Sprintf("PID %v", os.Getpid()))
    p.OpenStreams()
}

func (p *ChunkStream) OpenStreams() {
    p.f1 = p.OpenInputStreamReadOnly(p.Src)
    p.f2 = p.OpenOutputStreamReadOnly(p.Dst)
    p.f3 = p.OpenOutputStreamReadWrite(p.Dst)
}

func (p *ChunkStream) OpenInputStreamReadOnly(s string) (*os.File) {

    if ! p.ExistsStream(s) {
        p.Die(fmt.Sprintf("Input %s does not exist ", s))
    }
    f, e := os.Open(s)
    p.DieOnError(e, fmt.Sprintf("Opening input ro %v ", s))
    return f
}

func (p *ChunkStream) OpenOutputStreamReadOnly(s string) (*os.File) {
    if ! p.ExistsStream(s) {
        if ! p.Force {
            p.Die(fmt.Sprintf("Output %s does not exist. Use -force to write file. ", s))
        }
        p.CreateStream(s)
    }
    f, e := os.Open(s)
    p.DieOnError(e, fmt.Sprintf("Opening output ro %v ", s))
    return f
}

func (p *ChunkStream) OpenOutputStreamReadWrite(s string)  (*os.File) {
    f, e := os.OpenFile(s, os.O_RDWR, 0644)
    p.DieOnError(e, fmt.Sprintf("Opening output rw %v ", s))
    return f
}

func (p *ChunkStream) ExistsStream(path string) (exists bool) {
  exists = true
  if _, err := os.Stat(path); os.IsNotExist(err) {
    exists = false
  } 
  return
}

func (p *ChunkStream) CreateStream(s string) {
    var f, e = os.Create(s)
    p.DieOnError(e, fmt.Sprintf("Creating %v ", s))
    defer f.Close()
}

func (p *ChunkStream) SyncStream(f *os.File, s string) {
    if p.Flush {
        p.DieOnError(f.Sync(), fmt.Sprintf("Syncing %v ", s))
    }
}

func (p *ChunkStream) CloseStream(f *os.File, s string) {
    p.DieOnError(f.Close(), fmt.Sprintf("Closing %v ", s))
}

func (p *ChunkStream) CloseStreams() {
    p.CloseStream(p.f1, p.Src)
    p.CloseStream(p.f2, p.Dst)
    p.SyncStream(p.f3, p.Dst)
    p.CloseStream(p.f3, p.Dst)
}

func (p *ChunkStream) Die(msg string) {
    p.Println(msg)
    p.CloseStreams()
    p.Statistics()
    os.Exit(0) // TODO return appropriate values 
}

func (p *ChunkStream) DieOnError(e error, s string) {
    if e != nil {
        fmt.Sprintf("%v failed with error %v ", s, e)
        os.Exit(0) // TODO return appropriate values
    } else {
        p.Log("debug", s)
    } 
}

func (p *ChunkStream) DieOnSignal(msg string) {
    if p.Debug {
        p.Log("debug", msg)
    } else {
        p.Println("")
    }
    p.Statistics()
    os.Exit(0) // TODO return appropriate values 
}

func (p *ChunkStream) Log(level string, msg string) {
    if level == "verbose" {
        if(p.Verbose) {
            p.Println(msg)
        }
    }
    if level == "debug" {
        if(p.Debug) {
            p.Println(msg)
        }
    }
}

func (p *ChunkStream) Println(msg string) {
    if p.Uselog {
        log.Println(msg)
    } else {
        fmt.Println(msg)
    }
}

func (p *ChunkStream) ReadSourceChunk(b []byte, f *os.File) int {

    n, err := f.Read(b)

    if n == 0 {
        if err == nil {
        }
        if err == io.EOF {
            p.eofOnSrc = true
        } else {
            p.Die(fmt.Sprintf("Error while reading input chunk: %v", err))
        }
    }

    if ! p.eofOnSrc {
        p.bytesIn += int64(n)  
        p.chunksIn++
    }

    return n
}

func (p *ChunkStream) ReadDestinationChunk(b []byte, f *os.File) int {

    n, err := f.Read(b)

    if n == 0 {
        if err == nil {
        }
        if err == io.EOF {
            p.eofOnDest = true
        } else { 
            p.Die(fmt.Sprintf("Error while reading output chunk: %v", err))
        }
    } 

    return n
}

func (p *ChunkStream) HashChunk(b []byte) ([]byte) {
    hash := sha256.New()
    hash.Reset()
    hash.Write(b)
    return hash.Sum(nil)
}

func (p *ChunkStream) WriteChunk(b1 []byte, n1 int, offsetOut int64, f *os.File) {

    n2 := p.Write(b1, n1, offsetOut, f)

    if n1 != n2 {
        p.Die(fmt.Sprintf("Bytes read from src=%d and writen to dst=%d not equal", n1, n2))
    }
}

func (p *ChunkStream) Write(b []byte, i int, offsetOut int64, f *os.File) int {

    if p.Simulate {
        return 0    
    }

    n, err := f.Write(b[:i])

    if err != nil {  
        p.Die(fmt.Sprintf("Error while writing: %v", err))
    }
    if n != i {
        p.Die(fmt.Sprintf("Can't write chunk on dst=%v with chunksize from src=%v", n, i))
    }

    p.bytesOut += int64(n)  
    p.chunksOut++

    return n
}

func (p *ChunkStream) Seek(offsetIn int64, f *os.File) int64 {

    offsetOut, err := f.Seek(offsetIn, 0)

    if err != nil {  
        p.Die(fmt.Sprintf("Error while seeking: %v", err))
    }
    if !(offsetIn == offsetOut) {
        p.Die(fmt.Sprintf("Can't seek on offset from src=%d on dst=%d", 
            offsetIn, offsetOut))
    }
    return offsetOut
}

func (p *ChunkStream) Round(f float64, places int) (float64) {
    shift := math.Pow(10, float64(places))
    return math.Floor((f * shift) + .5) / shift;    
}

func (p *ChunkStream) Statistics() {
    if p.params.Stats {
        timeElapsed := time.Since(p.timeStart)
        mbPerSeconds := (float64(p.bytesIn)/1000000)/(float64(timeElapsed)/float64(time.Second))
        mbPerSeconds = p.Round(mbPerSeconds, 2)
	logger.Log(logger.Notice, fmt.Sprintf("domain %s: %v Bytes readed, %v Bytes writen, %v Chunks readed, %v Chunks writen, %s, %v MB/s ", 
            p.params.Domain, p.bytesIn, p.bytesOut, p.chunksIn, p.chunksOut, timeElapsed, mbPerSeconds))
    }
}

func (p *ChunkStream) Execute() {

    defer p.CloseStream(p.f3, p.Dst)
    defer p.SyncStream(p.f3, p.Dst)
    defer p.CloseStream(p.f2, p.Dst)
    defer p.CloseStream(p.f1, p.Src)

    p.chunksIn = 0
    p.chunksOut = 0
    p.bytesIn = 0
    p.bytesOut = 0

    p.eofOnSrc = false
    p.eofOnDest = false

    offsetIn := int64(0)

    b1 := make([]byte, p.ChunkSize)
    b2 := make([]byte, p.ChunkSize)
 
    for {
        n1 := p.ReadSourceChunk(b1, p.f1)

        if p.eofOnSrc {
            p.Log("debug", fmt.Sprintf("EOF on %v", p.Src))
            break
        }

        if n1 != p.ChunkSize {
            p.Log("debug", fmt.Sprintf("Correcting input chunksize from %d to %d", p.ChunkSize, n1))
            p.ChunkSize = n1
        }

        n2 := p.ReadDestinationChunk(b2, p.f2)

        offsetOut := p.Seek(offsetIn, p.f3)

        if p.eofOnDest {
            if ! p.Force {
                p.Die(fmt.Sprintf("EOF on %v while reading %v. Use -force to overwrite %v ", p.Dst, p.Src, p.Dst))
            }
            if p.Hash {
                p.Log("debug", fmt.Sprintf("Writing: chunk %d, offset %d, size %d, hash %v", 
                    p.chunksIn, offsetIn, p.ChunkSize, base64.StdEncoding.EncodeToString(p.HashChunk(b1))))
            } else {
                p.Log("debug", fmt.Sprintf("Writing: chunk %d, offset %d, size %d", 
                    p.chunksIn, offsetOut, p.ChunkSize))
            }
            p.WriteChunk(b1, n1, offsetOut, p.f3) 
            offsetIn = offsetIn + int64(p.ChunkSize) 
            continue 
        } 

        if n1 != n2 {
            if ! p.Force {
                p.Die(fmt.Sprintf("Bytes read from src=%d and dst=%d not equal. Use -force to overwrite %v ", n1, n2, p.Dst))
            }
            if n2 != p.ChunkSize { 
                //p.Log("debug", fmt.Sprintf("Correcting destination chunksize from %d to %d", n2, n1))
                n2 = n1 // currently not used 
            }
        }

        if p.Hash {
            h1 := p.HashChunk(b1)
            h2 := p.HashChunk(b2)
            if ! (bytes.Equal(h1, h2)) {
                p.WriteChunk(b1, n1, offsetOut, p.f3)
                //p.Log("debug", fmt.Sprintf("Writing: chunk %d, offset %d, size %d, hash %v, hash %v", 
                    //p.chunksIn, offsetIn, p.ChunkSize, base64.StdEncoding.EncodeToString(h1), base64.StdEncoding.EncodeToString(h2)))
            } else {
                //p.Log("debug", fmt.Sprintf("Equal: chunk %d, offset %d, size %d, hash %v, hash %v", 
                    //p.chunksIn, offsetIn, p.ChunkSize, base64.StdEncoding.EncodeToString(h1), base64.StdEncoding.EncodeToString(h2)))
            } 
        } else {
            if ! bytes.Equal(b1, b2) {
                p.WriteChunk(b1, n1, offsetOut, p.f3)
                //p.Log("debug", fmt.Sprintf("Writing: chunk %d, offset %d, size %d", 
                    //p.chunksIn, offsetOut, p.ChunkSize))
            } else {
                //p.Log("debug", fmt.Sprintf("Equal: chunk %d, offset %d, size %d", 
                    //p.chunksIn, offsetIn, p.ChunkSize))
            }
        }
        offsetIn = offsetIn + int64(p.ChunkSize)
    }
}

func(p *ChunkStream) ProcessSignals(signal_chan chan os.Signal) {
    for {
        sigtype := <-signal_chan
        switch sigtype {
            case syscall.SIGUSR1:
                p.Log("debug", "\nSignal sigusr received")
                p.Statistics()
            case syscall.SIGHUP:
                p.DieOnSignal("\nSignal hungup received")
            case syscall.SIGINT:
                p.DieOnSignal("\nSignal sigint received")
            case syscall.SIGTERM:
                p.DieOnSignal("\nSignal sigterm received")
            case syscall.SIGQUIT:
                p.DieOnSignal("\nSignal sigquit received")
            default:
                p.DieOnSignal("\nUnknown signal received")
        }
    }
}
