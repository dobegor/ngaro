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

package vm

import (
	"github.com/pkg/errors"
)

// Ngaro Virtual Machine Opcodes.
const (
	OpNop Cell = iota
	OpLit
	OpDup
	OpDrop
	OpSwap
	OpPush
	OpPop
	OpLoop
	OpJump
	OpReturn
	OpGtJump
	OpLtJump
	OpNeJump
	OpEqJump
	OpFetch
	OpStore
	OpAdd
	OpSub
	OpMul
	OpDimod
	OpAnd
	OpOr
	OpXor
	OpShl
	OpShr
	OpZeroExit
	OpInc
	OpDec
	OpIn
	OpOut
	OpWait
	OpCall
	OpFAdd
	OpFSub
	OpFMul
	OpFDiv
	OpFtoi
	OpItof
	OpFGtJump
	OpFLtJump
	OpFNeJump
	OpFEqJump
)

// Tos returns the value of the Top item On the data Stack. Always returns 0 if
// Instance.Depth() is 0.
func (i *Instance) Tos() Cell {
	return i.tos
}

// SetTos sets (changes) the value of the Top item On the data Stack. If the stack
// is empty, this function will be a no-op (i.e. Tos() will return 0).
func (i *Instance) SetTos(v Cell) {
	if i.Depth() > 0 {
		i.tos = v
	}
}

// Nos returns the value of the Next item On the data Stack. Always returns 0 if
// Instance.Depth() is less than 2.
func (i *Instance) Nos() Cell {
	return i.data[i.sp]
}

// Depth returns the data stack depth.
func (i *Instance) Depth() int {
	return i.sp
}

// RDepth returns the address stack depth.
func (i *Instance) RDepth() int {
	return i.rsp
}

// Drop2 removes the top two items from the data stack.
func (i *Instance) Drop2() {
	i.sp -= 2
	if i.sp < 0 {
		panic(errors.New("data stack underflow"))
	}
	i.tos = i.data[i.sp+1] // NOTE: this works because i.data[0:2] is always 0
}

// Push pushes the argument on top of the data stack.
func (i *Instance) Push(v Cell) {
	i.sp++
	i.data[i.sp], i.tos = i.tos, v
}

// Pop pops the value on top of the data stack and returns it.
func (i *Instance) Pop() Cell {
	if i.sp == 0 {
		panic(errors.New("data stack underflow"))
	}

	tos := i.tos
	i.tos = i.data[i.sp]
	i.sp--
	return tos
}

// Rpush pushes the argument on top of the address stack.
func (i *Instance) Rpush(v Cell) {
	i.rsp++
	i.address[i.rsp], i.rtos = i.rtos, v
}

// Rpop pops the value on top of the address stack and returns it.
func (i *Instance) Rpop() Cell {
	if i.rsp == 0 {
		panic(errors.New("return stack underflow"))
	}

	rtos := i.rtos
	i.rtos = i.address[i.rsp]
	i.rsp--
	return rtos
}

func (i *Instance) Stop() <-chan struct{} {
	i.stopCh = make(chan struct{})
	return i.stopCh
}

func (i *Instance) Stopped() bool {
	return i.stopped
}

// Run starts execution of the VM.
//
// If an error occurs, the PC will will point to the instruction that triggered
// the error. The most likely error condition is an "index out of range"
// runtime.Error which can occur in the following cases:
//
//	- address or data stack full
//	- attempt to address memory outside of the range [0:len(i.Image)]
//	- use of a port number outside of the range [0:1024] in an I/O operation.
//
//	A full stack trace should be obtainable with:
//
//	fmt.Sprintf("%+v", err)
//
// The VM will not error on stack underflows. i.e. drop always succeeds, and
// both Instance.Tos() and Instance.Nos() on an empty stack always return 0. This
// is a design choice that enables end users to use the VM interactively with
// Retro without crashes on stack underflows.
//
// Please note that this behavior should not be used as a feature since it may
// change without notice in future releases.
//
// If the VM was exited cleanly from a user program with the `bye` word, the PC
// will be equal to len(i.Image) and err will be nil.
//
// Note that this package makes heavy use of the github.com/pkg/errors package.
// The "root cause" error can be obtained with errors.Cause().
//
// If the last input stream gets closed, the VM will exit and the root cause
// error will be io.EOF. This is a normal exit condition in most use cases.
func (i *Instance) Run() (err error) {
	i.stopped = false

	defer func() {
		if e := recover(); e != nil {
			switch e := e.(type) {
			case error:
				err = errors.Wrapf(e, "Recovered error @pc=%d/%d, stack %d/%d, rstack %d/%d",
					i.PC, len(i.Mem), i.sp, len(i.data)-1, i.rsp, len(i.address)-1)
			default:
				panic(e)
			}
		}
	}()

	i.insCount = 0
	for i.PC < len(i.Mem) {
		if i.stopCh != nil {
			i.stopped = true
			close(i.stopCh)
			i.stopCh = nil
			return nil
		}

		op := i.Mem[i.PC]
		switch op {
		case OpNop:
			i.PC++
		case OpLit:
			i.Push(i.Mem[i.PC+1])
			i.PC += 2
		case OpDup:
			i.sp++
			i.data[i.sp] = i.tos
			i.PC++
		case OpDrop:
			i.Pop()
			i.PC++
		case OpSwap:
			i.tos, i.data[i.sp] = i.data[i.sp], i.tos
			i.PC++
		case OpPush:
			i.Rpush(i.Pop())
			i.PC++
		case OpPop:
			i.Push(i.Rpop())
			i.PC++
		case OpLoop:
			v := i.tos - 1
			if v > 0 {
				i.tos = v
				i.PC = int(i.Mem[i.PC+1])
			} else {
				i.Pop()
				i.PC += 2
			}
		case OpJump:
			i.PC = int(i.Mem[i.PC+1])
		case OpReturn:
			i.PC = int(i.Rpop() + 1)
		case OpGtJump:
			if i.data[i.sp] > i.tos {
				i.PC = int(i.Mem[i.PC+1])
			} else {
				i.PC += 2
			}
			i.Drop2()
		case OpLtJump:
			if i.data[i.sp] < i.tos {
				i.PC = int(i.Mem[i.PC+1])
			} else {
				i.PC += 2
			}
			i.Drop2()
		case OpNeJump:
			if i.data[i.sp] != i.tos {
				i.PC = int(i.Mem[i.PC+1])
			} else {
				i.PC += 2
			}
			i.Drop2()
		case OpEqJump:
			if i.data[i.sp] == i.tos {
				i.PC = int(i.Mem[i.PC+1])
			} else {
				i.PC += 2
			}
			i.Drop2()
		case OpFetch:
			i.tos = i.Mem[i.tos]
			i.PC++
		case OpStore:
			i.Mem[i.tos] = i.data[i.sp]
			i.Drop2()
			i.PC++
		case OpAdd:
			rhs := i.Pop()
			i.tos += rhs
			i.PC++
		case OpSub:
			rhs := i.Pop()
			i.tos -= rhs
			i.PC++
		case OpMul:
			rhs := i.Pop()
			i.tos *= rhs
			i.PC++
		case OpDimod:
			lhs, rhs := i.data[i.sp], i.tos
			i.data[i.sp] = lhs % rhs
			i.tos = lhs / rhs
			i.PC++
		case OpAnd:
			rhs := i.Pop()
			i.tos &= rhs
			i.PC++
		case OpOr:
			rhs := i.Pop()
			i.tos |= rhs
			i.PC++
		case OpXor:
			rhs := i.Pop()
			i.tos ^= rhs
			i.PC++
		case OpShl:
			rhs := i.Pop()
			i.tos <<= uint8(rhs)
			i.PC++
		case OpShr:
			rhs := i.Pop()
			i.tos >>= uint8(rhs)
			i.PC++
		case OpZeroExit:
			if i.tos == 0 {
				i.PC = int(i.Rpop() + 1)
				i.Pop()
			} else {
				i.PC++
			}
		case OpInc:
			i.tos++
			i.PC++
		case OpDec:
			i.tos--
			i.PC++
		case OpIn:
			port := i.tos
			if h := i.inH[port]; h != nil {
				i.Pop()
				if err = h(i, port); err != nil {
					return errors.Wrap(err, "IN failed")
				}
			} else {
				// we're not calling i.In so that we can optimize out a Pop/Push
				// sequence
				i.tos, i.Ports[port] = i.Ports[port], 0
			}
			i.PC++
		case OpOut:
			v, port := i.data[i.sp], i.tos
			i.Drop2()
			if h := i.outH[port]; h != nil {
				err = h(i, v, port)
			} else {
				err = i.Out(v, port)
			}
			if err != nil {
				return errors.Wrap(err, "OUT failed")
			}
			i.PC++
		case OpWait:
			if i.Ports[0] != 1 {
				for p, h := range i.waitH {
					v := i.Ports[p]
					if v == 0 {
						continue
					}
					if err = h(i, v, p); err != nil {
						return errors.Wrap(err, "WAIT failed")
					}
				}
			}
			i.PC++

		// Extended opcodes
		case OpCall:
			i.Rpush(Cell(i.PC + 1))
			i.PC = int(i.Mem[i.PC+1])
		case OpFAdd:
			rhs := i.Pop()
			*i.tos.AsFCell() += *rhs.AsFCell()
			i.PC++
		case OpFSub:
			rhs := i.Pop()
			*i.tos.AsFCell() -= *rhs.AsFCell()
			i.PC++
		case OpFMul:
			rhs := i.Pop()
			*i.tos.AsFCell() *= *rhs.AsFCell()
			i.PC++
		case OpFDiv:
			rhs := i.Pop()
			*i.tos.AsFCell() /= *rhs.AsFCell()
			i.PC++
		case OpItof:
			*i.tos.AsFCell() = FCell(i.tos)
			i.PC++
		case OpFtoi:
			i.tos = Cell(*i.tos.AsFCell())
			i.PC++
		case OpFGtJump:
			if *i.data[i.sp].AsFCell() > *i.tos.AsFCell() {
				i.PC = int(i.Mem[i.PC+1])
			} else {
				i.PC += 2
			}
			i.Drop2()
		case OpFLtJump:
			if *i.data[i.sp].AsFCell() < *i.tos.AsFCell() {
				i.PC = int(i.Mem[i.PC+1])
			} else {
				i.PC += 2
			}
			i.Drop2()
		case OpFNeJump:
			if *i.data[i.sp].AsFCell() != *i.tos.AsFCell() {
				i.PC = int(i.Mem[i.PC+1])
			} else {
				i.PC += 2
			}
			i.Drop2()
		case OpFEqJump:
			if *i.data[i.sp].AsFCell() == *i.tos.AsFCell() {
				i.PC = int(i.Mem[i.PC+1])
			} else {
				i.PC += 2
			}
			i.Drop2()

		default:
			if op < 0 && i.opHandler != nil {
				// custom opcode
				err = i.opHandler(i, op)
				if err != nil {
					return errors.Wrap(err, "custom opcode handler failed")
				}
				i.PC++
			} else {
				return errors.Errorf("invalid opcode %d", op)
			}
		}
		i.insCount++
		if i.tickFn != nil && i.insCount&i.tickMask == 0 {
			i.tickFn(i)
		}
	}
	return nil
}
