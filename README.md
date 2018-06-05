[![Go Report Card](https://goreportcard.com/badge/github.com/dobegor/ngaro)](https://goreportcard.com/report/github.com/dobegor/ngaro)
[![GoDoc](https://godoc.org/github.com/dobegor/ngaro/vm?status.svg)](https://godoc.org/github.com/dobegor/ngaro/vm)

# Ngaro Go

## <a name="pkg-overview">Overview</a>
This is a fork of an original embeddable Go implementation of the [Ngaro Virtual Machine](http://retroforth.org/docs/The_Ngaro_Virtual_Machine.html) by [Denis Bernard](https://github.com/db47h/ngaro).

Fork was created to: 
- add float opcodes support into VM 
- add support to limit performance via `ClockPeriod` option
- remove implicit calls for cells with values over 30 
- adding an explicit `OpCall` opcode
- implementing `Stop()` method to gracefully pause the VM with ability to continue with `Run()`

Also the 32 bit support was phased out.

Current TODO:
- [x] make the symbolic assembler able to parse float literals
- [ ] interrupts support (`OpIRQ` and `OpINT` opcodes, external host-vm interrupts, runtime exception interrupts, such as divide by zero)
- [ ] ability to dump / restore full VM state (memory, stacks, instruction pointer)
- [ ] a simple language compiler with this VM as a target
- [ ] update assembler docs and VM spec accordingly

This repository contains the embeddable [virtual
machine](https://godoc.org/github.com/dobegor/ngaro/vm) and a rudimentary
[symbolic assembler](https://godoc.org/github.com/dobegor/ngaro/asm)
for easy bootstrapping of projects written in Ngaro machine language.

The main purpose of this implementation is to allow customization and
communication between embeddable programs and Go programs via custom opcodes and
I/O handlers. The package examples demonstrate various use cases. 
For more details on I/O handling in the Ngaro VM, please refer to http://retroforth.org/docs/The_Ngaro_Virtual_Machine.html.

Another goal is to make the VM core as neutral as possible regarding the higher
level language running on it. For example, the in-memory string encoding scheme
is fully customizable.

Custom opcodes are implemented by providing a opcode handler for cells with negative values.
The maximum number of addressable cells is 2^63. The range [-2^63 - 1, -1] is available
for custom opcodes.

For all intents and purposes, the VM behaves according to the specification, except the 
aforementioned changes in this fork.

This is of particular importance to implementors of custom opcodes: the VM
always increments the PC after each opcode, thus opcodes altering the PC must
adjust it accordingly (i.e. set it to its real target minus one).

## Installing

Install the library:

	go get -u -v github.com/dobegor/ngaro/...

Test:

	go test -i github.com/dobegor/ngaro/vm
	go test -v github.com/dobegor/ngaro/vm/...

## Contributing

No rules, just common sense. Bells, wristles and any other preformance improvements are very
welcome. The only changes that will never be accepted are those that will break compatibility
with the VM specification.

PRs are a good place to discuss changes so do not hesitate to send PRs directly. Or an issue, if your really have an issue.

## License

This project is Copyright 2016 Denis Bernard <db047h@gmail.com>, licensed under
the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).

The Ngaro Virtual Machine is Copyright (c) 2008-2016 Charles
Childers (and many others), licensed under the ISC license.

## Contributors
[Pavel Vasilev](https://github.com/exploser/)

[George Dobrovolsky](https://github.com/dobegor/)
