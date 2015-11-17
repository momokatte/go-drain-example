// Copyright 2015 Michael O'Rourke. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"sync"

	drain "github.com/momokatte/go-drain"
)

var outStrings chan string

// Example app which reads strings from STDIN or a file, processes them concurrently, and writes them to STDOUT
// or a file.
//
func main() {

	var inFile, outFile string
	flag.StringVar(&inFile, "infile", "", "Input filename")
	flag.StringVar(&outFile, "outfile", "", "Output filename")
	flag.Parse()

	var outWG sync.WaitGroup
	outStrings = make(chan string, 1024)
	outWG.Add(1)
	go func() {
		if len(outFile) > 0 {
			if err := drain.ChanToFile(outStrings, outFile); err != nil {
				panic(err)
			}
		} else {
			if err := drain.ChanToStdout(outStrings); err != nil {
				panic(err)
			}
		}
		outWG.Done()
	}()

	inStrings := make(chan string, 200)
	go func() {
		if len(inFile) > 0 {
			if err := drain.FileLinesToChan(inFile, inStrings); err != nil {
				panic(err)
			}
		} else {
			if err := drain.StdinLinesToChan(inStrings); err != nil {
				panic(err)
			}
		}
		close(inStrings)
	}()

	go func() {
		for s := range inStrings {
			// process each string here
			outStrings <- s
		}
		close(outStrings)
	}()

	outWG.Wait()
}
