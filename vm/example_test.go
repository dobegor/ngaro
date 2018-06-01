// This file is part of ngaro - https://github.com/dobegor/ngaro
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

package vm_test

import (
	"fmt"
	"strings"

	"github.com/dobegor/ngaro/asm"
	"github.com/dobegor/ngaro/vm"
)

// Demonstrates how to use custom opcodes. This example defines a custom opcode
// that pushes the n-th fibonacci number onto the stack.
func ExampleBindOpcodeHandler() {
	fib := func(v vm.Cell) vm.Cell {
		var v0, v1 vm.Cell = 0, 1
		for v > 1 {
			v0, v1 = v1, v0+v1
			v--
		}
		return v1
	}

	handler := func(i *vm.Instance, opcode vm.Cell) error {
		switch opcode {
		case -1:
			i.SetTos(fib(i.Tos()))
			return nil
		default:
			return fmt.Errorf("Unsupported opcode value %d", opcode)
		}
	}

	img, err := asm.Assemble("test_fib_opcode", strings.NewReader(`
		.opcode fib -1	( define instruction fib as opcode -1 )
		46 fib
		`))
	if err != nil {
		panic(err)
	}

	i, err := vm.New(img, "dummy", vm.BindOpcodeHandler(handler))
	if err != nil {
		panic(err)
	}

	err = i.Run()
	if err != nil {
		panic(err)
	}

	// So, what's Fib(46)?
	fmt.Println(i.Data())

	// Output:
	// [1836311903]
}
