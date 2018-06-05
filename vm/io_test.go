// // This file is part of ngaro - https://github.com/dobegor/ngaro
// //
// // Copyright 2016 Denis Bernard <db047h@gmail.com>
// //
// // Licensed under the Apache License, Version 2.0 (the "License");
// // you may not use this file except in compliance with the License.
// // You may obtain a copy of the License at
// //
// //     http://www.apache.org/licenses/LICENSE-2.0
// //
// // Unless required by applicable law or agreed to in writing, software
// // distributed under the License is distributed on an "AS IS" BASIS,
// // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// // See the License for the specific language governing permissions and
// // limitations under the License.

package vm_test

// import (
// 	"bytes"
// 	"io"
// 	"os"
// 	"strings"
// 	"testing"
// 	"unsafe"

// 	"github.com/dobegor/ngaro/asm"
// 	"github.com/dobegor/ngaro/vm"
// 	"github.com/pkg/errors"
// )

// func Test_io_Stacks(t *testing.T) {
// 	i, err := runAsmImage("-16 5 out 0 0 out wait 5 in -17 5 out 0 0 out wait 5 in", "io_Stacks",
// 		vm.DataSize(24), vm.AddressSize(42))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	assertEqualI(t, "io_Stacks", 42, int(i.Pop()))
// 	assertEqualI(t, "io_Stacks", 24, int(i.Pop()))
// }

// func Test_io_Caps(t *testing.T) {
// 	// TODO: although the VM should return a correct value for endianness,
// 	// the test will fail on BigEndian platforms
// 	i, err := runAsmImage(`jump start
// 		.org 32
// 		:io dup 3 ! out 0 0 out wait 3 @ in ;
// 		:start
// 			-6 5 call io ( rstack size should be 1 (inside :io) )
// 			42 push 42 push -6 5 call io ( rstack size should be 3 (+1 inside :io) )
// 			-13 5 call io ( cell bits )
// 			-14 5 call io ( endianness )
// 			-15 5 call io ( port 8 enabled )
// 			 1 1 call io ( will cause EOF on nil input )`,
// 		"io_Caps", vm.Output(vm.NewVT100Terminal(bytes.NewBuffer(nil), nil, nil)))
// 	if errors.Cause(err) != io.EOF {
// 		t.Fatalf("Unexpected error: %v", err)
// 	}
// 	assertEqualI(t, "io_Caps port 8", -1, int(i.Pop()))
// 	assertEqualI(t, "io_Caps endian", 0, int(i.Pop()))
// 	assertEqualI(t, "io_Caps Cell bits", 8*int(unsafe.Sizeof(vm.Cell(0))), int(i.Pop())) // do not use vm.CellBits, just to check
// 	assertEqualI(t, "io_Caps rstack", 3, int(i.Pop()))
// 	assertEqualI(t, "io_Caps rstack", 1, int(i.Pop()))
// }

// // Test default In handler (not actually used in core for perf reasons).
// func TestVM_In(t *testing.T) {
// 	i, err := runAsmImage(`20 in 42 20 out 20 in 20 in`,
// 		"VM_In", vm.BindInHandler(20, (*vm.Instance).In))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	assertEqualI(t, "VM_In", 0, int(i.Pop()))
// 	assertEqualI(t, "VM_In", 42, int(i.Pop()))
// 	assertEqualI(t, "VM_In", 0, int(i.Pop()))
// }

// func Test_multireader(t *testing.T) {
// 	i, err := runAsmImage(`jump start
// 		.org 32
// 		:io dup push out 0 0 out wait pop in ;
// 		:start
// 		1 1 call io ( read from input until EOF )
// 		jump start`,
// 		"multireader",
// 		vm.Input(strings.NewReader("56")),
// 		vm.Input(strings.NewReader("34")),
// 		vm.Input(strings.NewReader("12")))
// 	if errors.Cause(err) != io.EOF {
// 		t.Fatalf("Unexpected error: %v", err)
// 	}
// 	for n := 6; n > 0; n-- {
// 		assertEqualI(t, "io_multireader", n+48, int(i.Pop()))
// 	}
// }

// func Test_port8(t *testing.T) {
// 	var flushed bool
// 	flush := func() error {
// 		flushed = true
// 		return nil
// 	}
// 	size := func() (int, int) { return 42, 24 }
// 	i, err := runAsmImage(`jump start
// 		.org 32
// 		:io dup push out 0 0 out wait pop in ;
// 		:start
// 		1 3 call io drop ( flush )
// 		-11 5 call io
// 		-12 5 call io
// 		-1 1 2 call io drop
// 		0 0 1 8 call io drop
// 		0 2 8 call io drop
// 		0 3 8 call io drop
// 		`,
// 		"port8",
// 		vm.Output(vm.NewVT100Terminal(bytes.NewBuffer(nil), flush, size)))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if !flushed {
// 		t.Fatal("Flush failed")
// 	}
// 	if i.Tos() != 24 {
// 		t.Fatalf("Expected height: 24, got: %d", i.Tos())
// 	}
// 	if i.Nos() != 42 {
// 		t.Fatalf("Expected width: 42, got: %d", i.Nos())
// 	}
// }

// func TestLoad(t *testing.T) {
// 	fn := "testdata/testLoad"
// 	f, err := os.Create(fn)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer os.Remove(fn)
// 	_, err = f.Write([]byte{0xff, 0xff, 0xff, 0xff, 0x01, 0x00, 0x00, 0x00})
// 	f.Close()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	mem, _, err := vm.Load(fn, 0, 32)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if mem[0] != vm.Cell(-1) {
// 		t.Fatalf("Expected -1, got %d", mem[0])
// 	}
// 	// force failure if vm.Cell is 32 bits
// 	if vm.CellBits == 32 {
// 		_, _, err = vm.Load(fn, 0, 64)
// 		exp := "load error: 64 bits value 8589934591 at memory location 0 too large"
// 		if err == nil || err.Error() != exp {
// 			t.Fatal(err)
// 		}
// 	}
// }

// func TestSave_64(t *testing.T) {
// 	d := "testdata/testDump64"
// 	img, err := asm.Assemble("Save", strings.NewReader("1 4 out 0 0 out wait 4 in"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer os.Remove(d)
// 	sf := func(fileName string, mem []vm.Cell) error {
// 		return vm.Save(fileName, mem, 64)
// 	}
// 	i, err := runImage(img, d, vm.SaveMemImage(sf))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	assertEqualI(t, "Save", 0, int(i.Pop()))
// 	saved, cells, err := vm.Load(d, 0, 64)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var same = true
// 	for n := range img {
// 		if img[n] != saved[n] {
// 			same = false
// 			break
// 		}
// 	}
// 	if !same {
// 		t.Fatalf("Save image error:\nexpected %v, got %v", img, saved[:cells])
// 	}
// }

// func TestSave_32(t *testing.T) {
// 	d := "testdata/testDump32"
// 	img, err := asm.Assemble("Save", strings.NewReader("jump start .dat 0 :start 1 4 out 0 0 out wait 4 in"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer os.Remove(d)
// 	sf := func(fileName string, mem []vm.Cell) error {
// 		return vm.Save(fileName, mem, 32)
// 	}
// 	// force failure if vm.Cell is 64 bits
// 	if vm.CellBits == 64 {
// 		x := int64(1)
// 		img[2] = vm.Cell(x << 32)
// 		_, err := runImage(img, d, vm.SaveMemImage(sf))
// 		exp := "WAIT failed: image dump failed: 64 bits value 4294967296 at memory location 2 too large"
// 		if err == nil || err.Error() != exp {
// 			t.Fatalf("\nExpected: %s\nGot: %v", exp, err.Error())
// 		}
// 		img[2] = 0
// 	}
// 	i, err := runImage(img, d, vm.SaveMemImage(sf))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	assertEqualI(t, "Save", 0, int(i.Pop()))
// 	saved, cells, err := vm.Load(d, 0, 32)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var same = true
// 	for n := range img {
// 		if img[n] != saved[n] {
// 			same = false
// 			break
// 		}
// 	}
// 	if !same {
// 		t.Fatalf("Save image error:\nexpected %v, got %v", img, saved[:cells])
// 	}
// }
