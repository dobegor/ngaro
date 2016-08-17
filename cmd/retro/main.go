// This file is part of ngaro - https://github.com/db47h/ngaro
//
// Copyright 2016 Denis Bernard <db047h@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/db47h/ngaro/vm"
)

var fileName = flag.String("image", "retroImage", "Load image from file `filename`")
var outFileName = flag.String("o", "", "Save image to file `filename`")
var withFile = flag.String("with", "", "Add `filename` to the input stack")
var shrink = flag.Bool("shrink", true, "When saving, don't save unused cells")
var size = flag.Int("size", 100000, "image size in cells")
var rawIO = flag.Bool("raw", true, "enable raw terminal IO")
var debug = flag.Bool("debug", false, "enable debug diagnostics")
var dump = flag.Bool("dump", false, "dump stacks and image upon exit, for ngarotest.py")

func port1Handler(i *vm.Instance, v, port vm.Cell) error {
	if v != 1 {
		return i.Wait(v, port)
	}
	// if v == 1, this will always read something
	e := i.Wait(v, port)
	// in raw tty mode, we need to handle CTRL-D ourselves
	if e == nil && i.Ports[1] == 4 {
		return io.EOF
	}
	return e
}

func port2Handler(w io.Writer) func(i *vm.Instance, v, port vm.Cell) error {
	return func(i *vm.Instance, v, port vm.Cell) error {
		var e error
		if v != 1 {
			return i.Wait(v, port)
		}
		t := i.Tos // save TOS
		e = i.Wait(v, port)
		if e == nil && t == 8 && i.Ports[port] == 0 {
			// the vm has written a backspace, erase char under cursor
			_, e = w.Write([]byte{32, 8})
		}
		return e
	}
}

func shrinkSave(mem []vm.Cell, fileName string) error {
	end := vm.Cell(len(mem))
	if len(mem) < 4 {
		return nil
	}
	if here := mem[3]; *shrink && here >= 0 && here < end {
		end = here
	}
	err := vm.Save(mem[:end], *outFileName)
	if err != nil {
		return err
	}
	return nil
}

func setupIO() (raw bool, tearDown func()) {
	var err error
	if *rawIO {
		tearDown, err = setRawIO()
		if err != nil {
			return false, nil
		}
	}
	return true, tearDown
}

func atExit(i *vm.Instance, err error) {
	if err == nil {
		return
	}
	if !*debug {
		fmt.Fprintf(os.Stderr, "\n%v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "\n%+v\n", err)
	if i != nil {
		if i.PC < len(i.Mem) {
			fmt.Fprintf(os.Stderr, "PC: %v (%v), Stack: %v, Addr: %v\n", i.PC, i.Mem[i.PC], i.Data(), i.Address())
		} else {
			fmt.Fprintf(os.Stderr, "PC: %v, Stack: %v\nAddr:  %v\n", i.PC, i.Data(), i.Address())
		}
	}
	os.Exit(1)
}

func main() {
	// check exit condition
	var err error
	var i *vm.Instance

	stdout := bufio.NewWriter(os.Stdout)
	output := vm.NewVT100Terminal(stdout, stdout.Flush, consoleSize(os.Stdout))

	// flush output, catch and log errors
	defer func() {
		output.Flush()
		atExit(i, err)
	}()

	flag.Parse()

	// try to switch the output terminal to raw mode.
	rawtty, ioTearDownFn := setupIO()
	if ioTearDownFn != nil {
		defer ioTearDownFn()
	}

	// default options
	var opts = []vm.Option{
		vm.SaveMemImage(shrinkSave),
		vm.Output(output),
	}

	if rawtty {
		// with the terminal in raw mode, we need to manually handle CTRL-D and
		// backspace, so we'll intercept WAITs on ports 1 and 2.
		// we could also do it with wrappers around Stdin/Stdout
		opts = append(opts,
			vm.Input(os.Stdin),
			vm.BindWaitHandler(1, port1Handler),
			vm.BindWaitHandler(2, port2Handler(output)))
	} else {
		// If not raw tty, buffer stdin, but do not check further if the i/o is
		// a terminal or not. The standard VT100 behavior is sufficient here.
		opts = append(opts, vm.Input(bufio.NewReader(os.Stdin)))
	}

	// append withFile to the input stack
	if len(*withFile) > 0 {
		var f *os.File
		f, err = os.Open(*withFile)
		if err != nil {
			return
		}
		opts = append(opts, vm.Input(bufio.NewReader(f)))
	}

	img, fileCells, err := vm.Load(*fileName, *size)
	if err != nil {
		return
	}
	if *outFileName == "" {
		outFileName = fileName
	}
	i, err = vm.New(img, *outFileName, opts...)
	if err != nil {
		return
	}
	if err = i.Run(); err == io.EOF {
		err = nil
	}
	if *dump {
		err = dumpVM(i, fileCells, output)
		if err != nil {
			return
		}
	}
}
