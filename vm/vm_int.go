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

import "unsafe"

// Cell is the basic type stored in a VM memory location
type Cell int

// FCell is the floating point representation of value stored in a VM memory location
type FCell float64

// AsFCell reinterprets Cell to FCell
func (c *Cell) AsFCell() *FCell {
	return (*FCell)(unsafe.Pointer(c))
}

// AsCell reinterprets FCell to Cell
func (f *FCell) AsCell() *Cell {
	return (*Cell)(unsafe.Pointer(f))
}
